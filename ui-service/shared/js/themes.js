// === BAR RESTAURANT - THEME CONFIGURATION ===

const THEMES = {
    // Classic Brown (Original)
    brown: {
        name: 'Classic Brown',
        primary: '#8B4513',
        primaryDark: '#6B3410',
        secondary: '#D4A574',
        accent: '#C9A227',
        accentSecondary: '#8B0000',
        dark: '#2C1810',
        backgroundLight: '#FDF5E6',
        icon: 'fa-wine-glass-alt'
    },
    
    // Navy & Gold - Classic, Sophisticated
    navy: {
        name: 'Navy & Gold',
        primary: '#1a2744',
        primaryDark: '#0f1829',
        secondary: '#3d5a80',
        accent: '#c9a227',
        accentSecondary: '#8b6914',
        dark: '#0a1628',
        backgroundLight: '#f0f4f8',
        icon: 'fa-martini-glass-citrus'
    },
    
    // Slate & Teal - Modern, Clean
    slate: {
        name: 'Slate & Teal',
        primary: '#2c3e50',
        primaryDark: '#1a252f',
        secondary: '#5d7b93',
        accent: '#00a8a8',
        accentSecondary: '#007777',
        dark: '#1a252f',
        backgroundLight: '#ecf0f1',
        icon: 'fa-champagne-glasses'
    },
    
    // Charcoal & Rose - Luxurious
    charcoal: {
        name: 'Charcoal & Rose',
        primary: '#2d2d2d',
        primaryDark: '#1a1a1a',
        secondary: '#4a4a4a',
        accent: '#b76e79',
        accentSecondary: '#8b4f57',
        dark: '#1a1a1a',
        backgroundLight: '#f5f0f0',
        icon: 'fa-wine-bottle'
    },
    
    // Deep Emerald - Rich, Upscale
    emerald: {
        name: 'Deep Emerald',
        primary: '#0d4f3c',
        primaryDark: '#083328',
        secondary: '#1a7a5c',
        accent: '#d4af37',
        accentSecondary: '#a88a2a',
        dark: '#052419',
        backgroundLight: '#f0f7f4',
        icon: 'fa-leaf'
    },
    
    // Midnight Purple - Royal
    purple: {
        name: 'Midnight Purple',
        primary: '#2d1b4e',
        primaryDark: '#1a1030',
        secondary: '#4a3072',
        accent: '#c0c0c0',
        accentSecondary: '#9370db',
        dark: '#1a1030',
        backgroundLight: '#f5f3f7',
        icon: 'fa-crown'
    },
    
    // Ocean Blue - Fresh, Professional
    ocean: {
        name: 'Ocean Blue',
        primary: '#0077b6',
        primaryDark: '#005a8c',
        secondary: '#0096c7',
        accent: '#ff9f1c',
        accentSecondary: '#f77f00',
        dark: '#023e5f',
        backgroundLight: '#f0f8ff',
        icon: 'fa-anchor'
    },
    
    // Burgundy - Wine Bar
    burgundy: {
        name: 'Burgundy Wine',
        primary: '#722f37',
        primaryDark: '#4a1f24',
        secondary: '#9b4d54',
        accent: '#d4af37',
        accentSecondary: '#c9a227',
        dark: '#2d1216',
        backgroundLight: '#faf5f5',
        icon: 'fa-wine-glass'
    }
};

// ============================================================================
// THEME MANAGER
// ============================================================================

class ThemeManager {
    constructor() {
        this.currentTheme = 'navy'; // Default theme
        this.storageKey = 'barrest_theme';
    }

    init() {
        // Load saved theme or use default
        const savedTheme = localStorage.getItem(this.storageKey);
        if (savedTheme && THEMES[savedTheme]) {
            this.currentTheme = savedTheme;
        }
        
        this.applyTheme(this.currentTheme);
        console.log(`🎨 Theme initialized: ${THEMES[this.currentTheme].name}`);
    }

    applyTheme(themeName) {
        const theme = THEMES[themeName];
        if (!theme) {
            console.error(`Theme "${themeName}" not found`);
            return;
        }

        const root = document.documentElement;
        
        // Apply CSS variables
        root.style.setProperty('--primary-color', theme.primary);
        root.style.setProperty('--primary-dark', theme.primaryDark);
        root.style.setProperty('--secondary-color', theme.secondary);
        root.style.setProperty('--accent-color', theme.accent);
        root.style.setProperty('--accent-secondary', theme.accentSecondary);
        root.style.setProperty('--dark-color', theme.dark);
        root.style.setProperty('--background-light', theme.backgroundLight);
        root.style.setProperty('--gradient-primary', `linear-gradient(135deg, ${theme.primary} 0%, ${theme.dark} 100%)`);
        root.style.setProperty('--gradient-accent', `linear-gradient(135deg, ${theme.accent} 0%, ${theme.accentSecondary} 100%)`);

        // Update icon if exists
        const brandIcon = document.querySelector('.brand-logo > i, .mobile-brand > i');
        if (brandIcon) {
            // Remove all fa- classes except 'fa' and size classes
            const classes = brandIcon.className.split(' ').filter(c => 
                c === 'fa' || c.startsWith('fa-') && (c.includes('x') || c === 'fas' || c === 'far' || c === 'fab')
            );
            brandIcon.className = `fas ${theme.icon} ${classes.filter(c => c.includes('x')).join(' ')}`;
        }

        this.currentTheme = themeName;
        localStorage.setItem(this.storageKey, themeName);
        
        console.log(`🎨 Theme applied: ${theme.name}`);
    }

    setTheme(themeName) {
        this.applyTheme(themeName);
    }

    getTheme() {
        return this.currentTheme;
    }

    getThemeInfo() {
        return THEMES[this.currentTheme];
    }

    getAllThemes() {
        return Object.keys(THEMES).map(key => ({
            id: key,
            ...THEMES[key]
        }));
    }

    // Create theme selector dropdown
    createThemeSelector(containerId) {
        const container = document.getElementById(containerId);
        if (!container) return;

        const select = document.createElement('select');
        select.className = 'form-select form-select-sm theme-selector';
        select.id = 'themeSelector';
        select.innerHTML = Object.keys(THEMES).map(key => 
            `<option value="${key}" ${key === this.currentTheme ? 'selected' : ''}>${THEMES[key].name}</option>`
        ).join('');

        select.addEventListener('change', (e) => {
            this.setTheme(e.target.value);
        });

        container.appendChild(select);
    }
}

// Global instance
const themeManager = new ThemeManager();

// Auto-initialize on DOM load
document.addEventListener('DOMContentLoaded', () => {
    themeManager.init();
});

// Export for global access
window.THEMES = THEMES;
window.themeManager = themeManager;

