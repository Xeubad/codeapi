package repository

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
	dbaccount "github.com/Wei-Shaw/sub2api/ent/account"
	dbaccountgroup "github.com/Wei-Shaw/sub2api/ent/accountgroup"
	dbgroup "github.com/Wei-Shaw/sub2api/ent/group"
	dbupstreamvendor "github.com/Wei-Shaw/sub2api/ent/upstreamvendor"
	"github.com/Wei-Shaw/sub2api/internal/domain"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type upstreamVendorRepository struct {
	client *ent.Client
}

func NewUpstreamVendorRepository(client *ent.Client) service.UpstreamVendorRepository {
	return &upstreamVendorRepository{client: client}
}

func (r *upstreamVendorRepository) List(ctx context.Context) ([]service.UpstreamVendor, error) {
	rows, err := r.client.UpstreamVendor.Query().
		Where(dbupstreamvendor.DeletedAtIsNil()).
		Order(ent.Asc(dbupstreamvendor.FieldName)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]service.UpstreamVendor, 0, len(rows))
	for _, row := range rows {
		out = append(out, toServiceUpstreamVendor(row))
	}
	return out, nil
}

func (r *upstreamVendorRepository) GetByID(ctx context.Context, id int64) (*service.UpstreamVendor, error) {
	row, err := r.client.UpstreamVendor.Query().
		Where(dbupstreamvendor.IDEQ(id), dbupstreamvendor.DeletedAtIsNil()).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, service.ErrUpstreamVendorNotFound
		}
		return nil, err
	}
	vendor := toServiceUpstreamVendor(row)
	return &vendor, nil
}

func (r *upstreamVendorRepository) Create(ctx context.Context, name string, budget int, notes *string) (*service.UpstreamVendor, error) {
	builder := r.client.UpstreamVendor.Create().
		SetName(name).
		SetConcurrencyBudget(budget)
	if notes != nil {
		builder.SetNotes(*notes)
	}
	row, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	vendor := toServiceUpstreamVendor(row)
	return &vendor, nil
}

func (r *upstreamVendorRepository) Update(ctx context.Context, id int64, name string, budget int, notes *string) (*service.UpstreamVendor, error) {
	builder := r.client.UpstreamVendor.UpdateOneID(id).
		SetName(name).
		SetConcurrencyBudget(budget)
	if notes != nil {
		builder.SetNotes(*notes)
	} else {
		builder.ClearNotes()
	}
	row, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, service.ErrUpstreamVendorNotFound
		}
		return nil, err
	}
	vendor := toServiceUpstreamVendor(row)
	return &vendor, nil
}

func (r *upstreamVendorRepository) SoftDelete(ctx context.Context, id int64) error {
	affected, err := r.client.UpstreamVendor.Update().
		Where(dbupstreamvendor.IDEQ(id), dbupstreamvendor.DeletedAtIsNil()).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return service.ErrUpstreamVendorNotFound
	}
	return nil
}

func (r *upstreamVendorRepository) ListEntries(ctx context.Context, vendorID int64) ([]service.UpstreamEntry, error) {
	rows, err := r.client.Account.Query().
		Where(
			dbaccount.TypeEQ(domain.AccountTypeUpstream),
			dbaccount.DeletedAtIsNil(),
		).
		Order(ent.Asc(dbaccount.FieldName)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	entries := make([]service.UpstreamEntry, 0, len(rows))
	accountIDs := make([]int64, 0, len(rows))
	byAccountID := make(map[int64]int, len(rows))
	for _, row := range rows {
		owner, ok := repoUpstreamVendorIDFromExtra(row.Extra)
		if !ok || owner != vendorID {
			continue
		}
		entry := service.UpstreamEntry{
			Groups:      []service.UpstreamEntryGroupRef{},
			AccountID:   row.ID,
			VendorID:    owner,
			Name:        row.Name,
			Platform:    row.Platform,
			BaseURL:     credentialString(row.Credentials, "base_url"),
			HasAPIKey:   credentialString(row.Credentials, "api_key") != "",
			Concurrency: row.Concurrency,
			Priority:    row.Priority,
			Schedulable: row.Schedulable,
			Status:      row.Status,
		}
		byAccountID[row.ID] = len(entries)
		entries = append(entries, entry)
		accountIDs = append(accountIDs, row.ID)
	}
	if len(accountIDs) == 0 {
		return entries, nil
	}
	links, err := r.client.AccountGroup.Query().
		Where(dbaccountgroup.AccountIDIn(accountIDs...)).
		WithGroup().
		Order(ent.Asc(dbaccountgroup.FieldGroupID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		idx, ok := byAccountID[link.AccountID]
		if !ok || link.Edges.Group == nil {
			continue
		}
		entries[idx].Groups = append(entries[idx].Groups, service.UpstreamEntryGroupRef{
			GroupID:   link.GroupID,
			GroupName: link.Edges.Group.Name,
			Platform:  link.Edges.Group.Platform,
			Priority:  link.Priority,
			Enabled:   link.Enabled,
		})
	}
	return entries, nil
}

func (r *upstreamVendorRepository) ListGroups(ctx context.Context) ([]service.UpstreamGroupOption, error) {
	rows, err := r.client.Group.Query().
		Where(dbgroup.DeletedAtIsNil()).
		Order(ent.Asc(dbgroup.FieldName)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]service.UpstreamGroupOption, 0, len(rows))
	for _, row := range rows {
		out = append(out, service.UpstreamGroupOption{
			ID:       row.ID,
			Name:     row.Name,
			Platform: row.Platform,
		})
	}
	return out, nil
}

func (r *upstreamVendorRepository) SetEntryGroupEnabled(ctx context.Context, accountID, groupID int64, enabled bool) error {
	affected, err := r.client.AccountGroup.Update().
		Where(
			dbaccountgroup.AccountIDEQ(accountID),
			dbaccountgroup.GroupIDEQ(groupID),
		).
		SetEnabled(enabled).
		Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return service.ErrUpstreamEntryNotFound
	}
	return nil
}

func (r *upstreamVendorRepository) SetEntryGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	if _, err := tx.AccountGroup.Delete().
		Where(dbaccountgroup.AccountIDEQ(accountID)).
		Exec(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if len(groupIDs) > 0 {
		builders := make([]*ent.AccountGroupCreate, 0, len(groupIDs))
		for _, groupID := range groupIDs {
			builders = append(builders, tx.AccountGroup.Create().
				SetAccountID(accountID).
				SetGroupID(groupID))
		}
		if err := tx.AccountGroup.CreateBulk(builders...).Exec(ctx); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func toServiceUpstreamVendor(row *ent.UpstreamVendor) service.UpstreamVendor {
	vendor := service.UpstreamVendor{
		ID:                row.ID,
		Name:              row.Name,
		ConcurrencyBudget: row.ConcurrencyBudget,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
	if row.Notes != nil {
		vendor.Notes = *row.Notes
	}
	return vendor
}

func credentialString(credentials map[string]any, key string) string {
	if credentials == nil {
		return ""
	}
	value, ok := credentials[key].(string)
	if !ok {
		return ""
	}
	return value
}

func repoUpstreamVendorIDFromExtra(extra map[string]any) (int64, bool) {
	if extra == nil {
		return 0, false
	}
	value, ok := extra[service.UpstreamVendorIDExtraKey]
	if !ok || value == nil {
		return 0, false
	}
	switch typed := value.(type) {
	case int64:
		return typed, typed > 0
	case int:
		return int64(typed), typed > 0
	case float64:
		id := int64(typed)
		return id, id > 0
	}
	return 0, false
}
