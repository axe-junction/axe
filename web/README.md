# React + TypeScript + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react) uses [Babel](https://babeljs.io/) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

# Axe Web - Modern Transport Planning App

A modern, responsive web application built with React.js and Tailwind CSS that provides comprehensive transport planning for Algiers, featuring interactive maps, route planning, and community features.

## 🚀 Features

### 🗺️ Interactive Map-Based Trip Planning

- **Leaflet-powered Maps**: Interactive maps focused on Algiers
- **Route Visualization**: Visual route rendering with start/end points
- **Transport Mode Filtering**: Filter by bus, metro, tram, and walking routes
- **Real-time Route Selection**: Click on routes to see detailed information

### 🔍 Smart Search & Filtering

- **Advanced Search**: Smart search with autocomplete for locations
- **Multiple Filters**: Filter by transport mode, duration, cost, and provider
- **Recent Searches**: Quick access to previously searched routes
- **Suggestions**: Dropdown suggestions for popular routes

### 🚌 Comprehensive Route Results

- **Detailed Trip Cards**: Provider, cost, duration, and stop information
- **Expandable Details**: Step-by-step journey breakdowns
- **Transfer Points**: Clear indication of where to change transport
- **Multiple Options**: Compare different route options

### 👤 User Profile Management

- **Profile Customization**: Edit name, email, phone, and location
- **Trip Statistics**: Track favorite routes and search history
- **Settings**: Manage notifications, privacy, and app preferences

### 🗣️ Community Features

- **Community Feed**: Share experiences, tips, and photos
- **Social Interactions**: Like and comment on posts
- **User Generated Content**: Real user reviews and recommendations
- **Travel Tips**: Community-driven transport advice

### ⭐ Favorites & History

- **Favorite Routes**: Save frequently used routes
- **Recent History**: Quick access to recent searches
- **Trip Management**: Organize and manage saved trips
- **Quick Stats**: Overview of usage patterns

## 🛠️ Tech Stack

- **Frontend**: React.js 19 with TypeScript
- **Styling**: Tailwind CSS v4 with custom design system
- **Routing**: React Router DOM v7
- **State Management**: Zustand for global state
- **Maps**: Leaflet with React-Leaflet for interactive maps
- **Icons**: Lucide React for consistent iconography
- **Build Tool**: Vite for fast development and building

## 🎨 Design System

- **Primary Color**: #6316DB (Purple)
- **Background**: #FFFFFF (White)
- **Typography**: Inter font family
- **Mobile-First**: Responsive design with mobile-first approach
- **Accessibility**: Focus states and keyboard navigation

## 📱 Responsive Design

### Mobile (< 768px)

- Bottom navigation bar
- Full-screen map view
- Collapsible search results
- Touch-optimized interactions

### Tablet (768px - 1024px)

- Adapted layout with side panels
- Enhanced touch targets
- Optimized for landscape/portrait

### Desktop (> 1024px)

- Side navigation
- Multi-panel layout
- Keyboard shortcuts
- Optimized for mouse interaction

## 🚀 Getting Started

### Prerequisites

- Node.js (v18 or higher)
- npm or yarn

### Installation

1. Clone the repository

```bash
git clone <repository-url>
cd axe-web
```

2. Install dependencies

```bash
npm install
```

3. Start the development server

```bash
npm run dev
```

4. Open [http://localhost:5173](http://localhost:5173) in your browser

### Building for Production

```bash
npm run build
```

### Running Tests

```bash
npm run test
```

### Linting

```bash
npm run lint
```

## 🗂️ Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── BottomNavigation.tsx
│   ├── Layout.tsx
│   └── TripCard.tsx
├── pages/              # Page components
│   ├── MapView.tsx
│   ├── SearchPage.tsx
│   ├── CommunityPage.tsx
│   ├── FavoritesPage.tsx
│   └── ProfilePage.tsx
├── store/              # State management
│   └── useAppStore.ts
├── assets/             # Static assets
├── App.tsx            # Main app component
├── main.tsx          # Entry point
└── index.css         # Global styles
```

## 🌟 Key Features Implementation

### State Management

- **Zustand Store**: Centralized state management for user data, trips, and app state
- **Persistence**: Local storage persistence for user preferences and favorites
- **Type Safety**: Full TypeScript support with proper typing

### Map Integration

- **Leaflet Maps**: Interactive maps with custom styling
- **Route Visualization**: Polyline rendering for routes
- **Markers**: Custom markers for start/end points
- **Focused on Algiers**: Optimized for Algiers city transport

### Component Architecture

- **Reusable Components**: Modular design with reusable UI components
- **Props Interface**: Well-defined TypeScript interfaces
- **Responsive Design**: Mobile-first responsive components

## 🎯 Future Enhancements

- **Real-time Data**: Integration with actual transport APIs
- **Offline Support**: Service worker for offline functionality
- **Push Notifications**: Real-time updates for route changes
- **Multi-language**: Support for Arabic and French
- **Dark Mode**: Theme switching capability
- **Advanced Analytics**: User behavior tracking and insights

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🌍 Algiers Transport Focus

This application is specifically designed for the public transport system in Algiers, Algeria, including:

- **ETUSA**: Bus network coverage
- **Metro**: Algiers Metro system
- **Tram**: Algiers tramway network
- **Walking Routes**: Pedestrian-friendly paths

Built with love for the people of Algiers 🇩🇿

You can also install [eslint-plugin-react-x](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-x) and [eslint-plugin-react-dom](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-dom) for React-specific lint rules:

```js
// eslint.config.js
import reactX from "eslint-plugin-react-x";
import reactDom from "eslint-plugin-react-dom";

export default tseslint.config([
  globalIgnores(["dist"]),
  {
    files: ["**/*.{ts,tsx}"],
    extends: [
      // Other configs...
      // Enable lint rules for React
      reactX.configs["recommended-typescript"],
      // Enable lint rules for React DOM
      reactDom.configs.recommended,
    ],
    languageOptions: {
      parserOptions: {
        project: ["./tsconfig.node.json", "./tsconfig.app.json"],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
]);
```
