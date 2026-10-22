import { createApp } from 'vue'
import App from './App.vue'
import { createPinia } from 'pinia'
import router from './router/index'
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
import { definePreset } from '@primevue/themes'
import ToastService from 'primevue/toastservice'

// NexaPreset: PrimeVue theme derived from Design_UI.md tokens
// Primary ramp built from --color-primary #1C3D5A at 600
// darkModeSelector: false because spec defines light mode only

const NexaPreset = definePreset(Aura, {
  semantic: { 
    primary: {
      50: '#F5F8FA',
      100: '#E8EEF3',
      200: '#D4DDE7',
      300: '#A3B9CC',
      400: '#5F82A3',
      500: '#3A5F7E',
      600: '#1C3D5A',
      700: '#162F44',
      800: '#10222F',
      900: '#0A151A'
    },
    colorScheme: {
      light: {
        primary: {
          color: '{primary.600}',
          contrastColor: '#FFFFFF',
          hoverColor: '{primary.500}',
          activeColor: '{primary.800}'
        },
        surface: {
          0: '#FFFFFF',
          50: '#F7F6F3',
          100: '#F7F6F3',
          200: '#E5E2DB',
          300: '#E5E2DB',
          400: '#E5E2DB',
          500: '#E5E2DB',
          600: '#E5E2DB',
          700: '#E5E2DB',
          800: '#E5E2DB',
          900: '#1F2937'
        },
        content: {
          background: '#FFFFFF',
          hoverBackground: '#F7F6F3',
          borderColor: '#E5E2DB',
          color: '#1F2937',
          hoverColor: '#1F2937'
        },
        text: {
          color: '#1F2937',
          hoverColor: '#0F2740',
          mutedColor: '#6B7280',
          hoverMutedColor: '#1F2937'
        },
        formField: {
          background: '#FFFFFF',
          disabledBackground: '#F7F6F3',
          filledBackground: '#F7F6F3',
          filledHoverBackground: '#F7F6F3',
          filledFocusBackground: '#FFFFFF',
          borderColor: '#E5E2DB',
          hoverBorderColor: '{primary.600}',
          focusBorderColor: '{primary.600}',
          invalidBorderColor: '#B54245',
          color: '#1F2937',
          disabledColor: '#9CA3AF',
          placeholderColor: '#9CA3AF',
          invalidPlaceholderColor: '#9CA3AF',
          floatLabelColor: '#6B7280',
          floatLabelFocusColor: '{primary.600}',
          floatLabelInvalidFocusColor: '#B54245',
          iconColor: '#6B7280',
          shadow: '0 0 0 0 transparent'
        },
        overlay: {
          modal: {
            background: '#FFFFFF',
            borderColor: '#E5E2DB',
            color: '#1F2937'
          }
        }
      }
    }
  },
  components: {
    button: {
      root: {
        borderRadius: '6px',
        paddingX: '12px',
        paddingY: '8px',
        fontSize: '14px',
        transitionDuration: '140ms'
      },
      sm: {
        fontSize: '12px',
        paddingX: '10px',
        paddingY: '6px'
      }
    },
    datatable: {
      root: {
        borderRadius: '10px',
        borderColor: '{content.border.color}',
        borderWidth: '1px'
      },
      header: {
        background: '{surface.0}',
        borderColor: '{content.border.color}',
        color: '{text.muted.color}',
        borderWidth: '0 0 1px 0',
        padding: '12px'
      },
      headerCell: {
        borderColor: '{content.border.color}',
        hoverColor: '{text.color}',
        selectedColor: '{text.color}',
        gap: '0.5rem',
        padding: '12px',
        focusRing: {
          width: '0',
          style: 'none',
          color: 'transparent',
          offset: '0'
        }
      },
      columnTitle: {
        fontWeight: '600'
      },
      row: {
        background: '{surface.0}',
        hoverBackground: '{surface.50}',
        selectedBackground: '{primary.50}',
        color: '{text.color}',
        hoverColor: '{text.color}',
        selectedColor: '{text.color}',
        focusRing: {
          width: '0',
          style: 'none',
          color: 'transparent',
          offset: '0'
        }
      },
      bodyCell: {
        borderColor: '{content.border.color}',
        padding: '12px',
        gap: '0.5rem'
      },
      footerCell: {
        background: '{surface.0}',
        borderColor: '{content.border.color}',
        color: '{text.color}',
        padding: '12px'
      },
      columnFooter: {
        fontWeight: '600'
      },
      footer: {
        background: '{surface.0}',
        borderColor: '{content.border.color}',
        color: '{text.color}',
        borderWidth: '1px 0 0 0',
        padding: '12px'
      }
    },
    dialog: {
      root: {
        background: '{surface.0}',
        borderColor: '{content.border.color}',
        color: '{text.color}',
        borderRadius: '14px',
        shadow: '0 12px 32px rgb(31 41 55 / 12%)',
        transitionDuration: '140ms'
      },
      header: {
        padding: '16px 20px',
        gap: '0.5rem'
      },
      title: {
        fontSize: '16px',
        fontWeight: '600'
      },
      content: {
        padding: '0 20px 20px 20px'
      },
      footer: {
        padding: '16px 20px',
        gap: '8px'
      }
    },
    select: {
      root: {
        borderRadius: '6px',
        transitionDuration: '140ms'
      },
      dropdown: {
        width: '2.5rem',
        color: '{text.muted.color}'
      }
    },
    autocomplete: {
      root: {
        transitionDuration: '140ms'
      }
    },
    datepicker: {
      root: {
        transitionDuration: '140ms'
      }
    }
  }
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(PrimeVue, {
  theme: {
    preset: NexaPreset,
    options: {
      darkModeSelector: false,
      cssLayer: false
    }
  },
  ripple: true
})
app.use(ToastService)
app.mount('#app')