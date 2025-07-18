import { Tabs } from 'expo-router';
import FontAwesome from '@expo/vector-icons/FontAwesome';

export default function AppLayout() {
  return (
    <Tabs
      screenOptions={{
        headerShown: false,
        tabBarStyle: {
          position: 'absolute',
          left: 10,
          right: 10,
          bottom: 16,
          backgroundColor: 'rgba(0,0,0,0.85)',
          borderRadius: 5,
          elevation: 8,
          borderTopWidth: 0,
          height: 70,
        },
        tabBarActiveTintColor: '#fff',
        tabBarInactiveTintColor: '#a0aec0',
        tabBarLabelStyle: {
          fontSize: 12,
          fontWeight: '500',
          marginTop: 2,
        },
        tabBarItemStyle: {
          marginHorizontal: 8,
          paddingVertical: 4,
          borderRadius: 8,
        },
        tabBarLabelPosition: 'below-icon',
      }}
    >
      <Tabs.Screen
        name="map/index"
        options={{
          title: 'map',
          tabBarIcon: ({ color }) => (
            <FontAwesome size={28} name="map-marker" color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="live/index"
        options={{
          title: 'Live',
          tabBarIcon: ({ color }) => (
            <FontAwesome size={28} name="compass" color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="profile/index"
        options={{
          title: 'Profile',
          tabBarIcon: ({ color }) => (
            <FontAwesome size={28} name="user" color={color} />
          ),
        }}
      />
    </Tabs>
  );
}