    import { Tabs } from 'expo-router';
    import FontAwesome from '@expo/vector-icons/FontAwesome';

    export default function AppLayout() {
    return (
        
        <Tabs screenOptions={{
        headerShown: false, 
        }}>

        <Tabs.Screen
            name="map/index.tsx"
            options={{
            title: 'Map', 
            tabBarIcon: ({ color }) => (
                <FontAwesome size={28} name="map-marker" color={color} />
            ),
            }}
        />
<Tabs.Screen
  name="Live/index.tsx"
  options={{
    title: 'Live',
    tabBarIcon: ({ color }) => (
      <FontAwesome size={28} name="compass" color={color} />
    ),
  }}
/>


        <Tabs.Screen
            name="profile/index.tsx"
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