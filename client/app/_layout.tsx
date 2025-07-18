import { Tabs } from "expo-router";
import { Ionicons } from '@expo/vector-icons';
import "../presentation/styles/global.css";

export default function RootLayout() {
  return (
    <Tabs
      screenOptions={{
        headerShown: false,
        tabBarStyle: {
          position: 'absolute',
          bottom: 0,
          left: 0,
          right: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.9)',
          borderTopWidth: 0,
          borderRadius: 0,
          height: 90,
          paddingBottom: 20,
          paddingTop: 10,
        },
        tabBarActiveTintColor: '#fff',
        tabBarInactiveTintColor: '#a0aec0',
        tabBarLabelStyle: {
          fontSize: 12,
          fontWeight: '500',
          marginTop: 4,
        },
      }}
    >
      <Tabs.Screen
        name="map/index"
        options={{
          title: 'Accueil',
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="home" size={size} color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="live/index"
        options={{
          title: 'En direct',
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="radio" size={size} color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="profile/index"
        options={{
          title: 'Profil',
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="person" size={size} color={color} />
          ),
        }}
      />
    </Tabs>
  );
}
