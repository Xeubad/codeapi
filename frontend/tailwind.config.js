/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // 强调色 - 深移青柠。对比度实测：500 对白字 3.19（原 teal-500 仅 2.49），
        // 600/700 落在 AA 4.5 以上，供强调文字使用。
        // 规则：只用于链接、选中条、图表首系列、焦点环；不作为大面积背景或主按钮底色。
        primary: {
          50: '#f6f9ec',
          100: '#eaf3d2',
          200: '#d5e9a8',
          300: '#b3d766',
          400: '#94c737',
          500: '#6f9e12',
          600: '#558310',
          700: '#4d7c0f',
          800: '#3f6212',
          900: '#2f4a0d',
          950: '#16220a'
        },
        // 中性族：亮色模式的 gray 与暗色模式的 dark/accent 原本分属纯灰与 slate 蓝灰两族，
        // 切换主题时色温会断层。这里统一成同一套冷灰墨黑。
        gray: {
          50: '#f6f7f9',
          100: '#eceef2',
          200: '#e7eaef',
          300: '#d8dce3',
          400: '#a8b0bd',
          500: '#7b8494',
          600: '#4a5261',
          700: '#2b323d',
          800: '#1b212b',
          900: '#111620',
          950: '#0b0e14'
        },
        // 辅助色与暗色背景共用同一族（历史上两者数值就完全相同）
        accent: {
          50: '#f6f7f9',
          100: '#eceef2',
          200: '#e7eaef',
          300: '#d8dce3',
          400: '#a8b0bd',
          500: '#7b8494',
          600: '#4a5261',
          700: '#2b323d',
          800: '#1b212b',
          900: '#111620',
          950: '#0b0e14'
        },
        dark: {
          50: '#f6f7f9',
          100: '#eceef2',
          200: '#e7eaef',
          300: '#d8dce3',
          400: '#a8b0bd',
          500: '#7b8494',
          600: '#4a5261',
          700: '#2b323d',
          800: '#1b212b',
          900: '#111620',
          950: '#0b0e14'
        }
      },
      fontFamily: {
        sans: [
          'system-ui',
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'Roboto',
          'Helvetica Neue',
          'Arial',
          'PingFang SC',
          'Hiragino Sans GB',
          'Microsoft YaHei',
          'sans-serif'
        ],
        mono: ['ui-monospace', 'SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', 'monospace']
      },
      boxShadow: {
        glass: '0 8px 32px rgba(0, 0, 0, 0.08)',
        'glass-sm': '0 4px 16px rgba(0, 0, 0, 0.06)',
        glow: '0 0 20px rgba(111, 158, 18, 0.22)',
        'glow-lg': '0 0 40px rgba(111, 158, 18, 0.32)',
        card: '0 1px 2px rgba(16, 24, 40, 0.05)',
        'card-hover': '0 4px 12px rgba(16, 24, 40, 0.08)',
        'inner-glow': 'inset 0 1px 0 rgba(255, 255, 255, 0.1)'
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-primary': 'linear-gradient(135deg, #6f9e12 0%, #4d7c0f 100%)',
        'gradient-dark': 'linear-gradient(135deg, #1b212b 0%, #0b0e14 100%)',
        'gradient-glass':
          'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
        'mesh-gradient':
          'radial-gradient(at 40% 20%, rgba(111, 158, 18, 0.10) 0px, transparent 50%), radial-gradient(at 80% 0%, rgba(77, 124, 15, 0.07) 0px, transparent 50%), radial-gradient(at 0% 50%, rgba(111, 158, 18, 0.06) 0px, transparent 50%)'
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'slide-down': 'slideDown 0.3s ease-out',
        'slide-in-right': 'slideInRight 0.3s ease-out',
        'scale-in': 'scaleIn 0.2s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        shimmer: 'shimmer 2s linear infinite',
        glow: 'glow 2s ease-in-out infinite alternate'
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' }
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideDown: {
          '0%': { opacity: '0', transform: 'translateY(-10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideInRight: {
          '0%': { opacity: '0', transform: 'translateX(20px)' },
          '100%': { opacity: '1', transform: 'translateX(0)' }
        },
        scaleIn: {
          '0%': { opacity: '0', transform: 'scale(0.95)' },
          '100%': { opacity: '1', transform: 'scale(1)' }
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' }
        },
        glow: {
          '0%': { boxShadow: '0 0 20px rgba(111, 158, 18, 0.25)' },
          '100%': { boxShadow: '0 0 30px rgba(111, 158, 18, 0.4)' }
        }
      },
      backdropBlur: {
        xs: '2px'
      },
      borderRadius: {
        lg: '7px',
        xl: '8px',
        '2xl': '10px',
        '4xl': '2rem'
      }
    }
  },
  plugins: []
}
