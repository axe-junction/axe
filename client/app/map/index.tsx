import React, { useState } from 'react';
import { View, Text, TouchableOpacity, SafeAreaView, TextInput } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import MapView, { Marker } from 'react-native-maps';
import { BlurView } from 'expo-blur';

export default function TransitApp() {
  const [currentScreen, setCurrentScreen] = useState<'map' | 'destination'>('map');
  const [fromLocation, setFromLocation] = useState('');
  const [toLocation, setToLocation] = useState('');

  const handleSearchClick = () => setCurrentScreen('destination');
  const handleBackToMap = () => setCurrentScreen('map');

  const recentDestinations = [
    { name: "1er Mai", location: "P2WH+RFW, Alger" },
    { name: "Sonatrach - Direction Générale", location: "P2WH+RFW, Hydra" },
  ];

  const transportModes = [
    { icon: 'train-outline', label: "Train", active: true },
    { icon: 'walk-outline', label: "Marche", active: false },
    { icon: 'bus-outline', label: "Tramway", active: true },
    { icon: 'subway-outline', label: "Métro", active: false },
    { icon: 'shield-checkmark-outline', label: "Mode sûr", active: true },
  ];

  const customMapStyle = [
    { elementType: "geometry", stylers: [{ color: "#EAEAEC" }] },
    { elementType: "labels.icon", stylers: [{ visibility: "off" }] },
    { elementType: "labels.text.fill", stylers: [{ color: "#9E9E9E" }] },
    { elementType: "labels.text.stroke", stylers: [{ color: "#FFFFFF" }] },
    { featureType: "administrative", elementType: "geometry", stylers: [{ color: "#ffffff" }] },
    { featureType: "poi", elementType: "geometry", stylers: [{ color: "#EAEAEC" }] },
    { featureType: "road", elementType: "geometry", stylers: [{ color: "#D6D6D6" }] },
    { featureType: "water", elementType: "geometry", stylers: [{ color: "#C4CDE2" }] },
  ];

  return (
    <View className="flex-1 relative">
      <MapView
        className="absolute inset-0"
        customMapStyle={customMapStyle}
        initialRegion={{
          latitude: 40.6828,
          longitude: -73.9754,
          latitudeDelta: 0.1,
          longitudeDelta: 0.1,
        }}
      >
        <Marker coordinate={{ latitude: 40.6828, longitude: -73.9754 }} pinColor="#6b46c1" />
        {currentScreen === 'map' && (
          <Marker coordinate={{ latitude: 40.6782, longitude: -73.9780 }} pinColor="#38a169" />
        )}
      </MapView>

      <SafeAreaView className="absolute top-0 inset-x-0 p-4 flex-row justify-between items-center">
        <Text className="text-base font-semibold text-black">9:41</Text>
        <View className="flex-row items-center">
          {[...Array(3)].map((_, i) => (
            <View key={i} className="w-[3px] h-[6px] rounded bg-black mr-[2px]" />
          ))}
          <View className="w-[3px] h-[6px] rounded bg-gray-400" />
        </View>
      </SafeAreaView>

      {currentScreen === 'map' ? (
        <>
          <BlurView intensity={10} className="absolute top-20 inset-x-5 rounded-3xl py-2 px-4 bg-white/70 flex-row items-center shadow-lg">
            <Ionicons name="wifi" size={16} color="#ed8936" />
            <Text className="ml-2 text-sm font-medium text-gray-600">
              vous êtes actuellement hors ligne
            </Text>
          </BlurView>

          <TouchableOpacity
            onPress={handleSearchClick}
            className="absolute top-1/2 inset-x-5 bg-purple-700 rounded-2xl py-4 items-center shadow-lg"
          >
            <Text className="text-white text-base font-semibold">
              Rechercher une destination
            </Text>
          </TouchableOpacity>

          <BlurView intensity={10} className="absolute bottom-24 inset-x-5 rounded-3xl p-5 bg-white/70 shadow-lg">
            <Text className="text-base font-semibold text-gray-900 mb-4">
              Destinations récentes
            </Text>
            {recentDestinations.map((d, i) => (
              <TouchableOpacity
                key={i}
                onPress={() => {
                  setToLocation(d.name);
                  handleSearchClick();
                }}
                className="flex-row items-center py-3"
              >
                <View className="w-10 h-10 rounded-full bg-purple-700 justify-center items-center mr-3">
                  <Ionicons name="location" size={20} color="#fff" />
                </View>
                <View className="flex-1">
                  <Text className="text-sm font-semibold text-gray-900">{d.name}</Text>
                  <Text className="text-xs text-gray-500 mt-0.5">{d.location}</Text>
                </View>
              </TouchableOpacity>
            ))}
          </BlurView>

          <BlurView intensity={20} className="absolute bottom-5 inset-x-5 rounded-3xl py-3 bg-black/80 flex-row justify-around items-center shadow-lg">
            <TouchableOpacity className="items-center">
              <Ionicons name="home" size={24} color="#fff" />
              <Text className="text-white text-xs font-medium mt-1">Accueil</Text>
            </TouchableOpacity>
            <TouchableOpacity className="items-center">
              <View className="w-6 h-6 rounded-full border-2 border-white justify-center items-center">
                <View className="w-2 h-2 rounded-full bg-white" />
              </View>
              <Text className="text-white text-xs font-medium mt-1">En direct</Text>
            </TouchableOpacity>
            <TouchableOpacity className="items-center">
              <Ionicons name="person" size={24} color="#a0aec0" />
              <Text className="text-gray-400 text-xs font-medium mt-1">Profil</Text>
            </TouchableOpacity>
          </BlurView>
        </>
      ) : (
        <BlurView intensity={10} className="absolute top-20 inset-x-5 rounded-3xl p-5 bg-white/70 shadow-lg">
          <View className="flex-row items-center mb-5">
            <TouchableOpacity onPress={handleBackToMap} className="mr-2">
              <Ionicons name="arrow-back" size={24} color="#6b46c1" />
            </TouchableOpacity>
            <Text className="text-lg font-bold text-gray-900">
              Choisissez votre destination
            </Text>
          </View>

          <View className="mb-5 space-y-3">
            <BlurView intensity={10} className="relative rounded-xl overflow-hidden">
              <View className="absolute left-3 top-4 w-1.5 h-1.5 rounded-full bg-gray-400 z-10" />
              <TextInput
                placeholder="Départ"
                placeholderTextColor="rgba(107,114,128,1)"
                className="py-3.5 pl-8 text-base text-gray-900 bg-gray-100/80"
                value={fromLocation}
                onChangeText={setFromLocation}
              />
            </BlurView>
            <BlurView intensity={10} className="rounded-xl overflow-hidden">
              <TextInput
                placeholder="Destination"
                placeholderTextColor="rgba(255,255,255,1)"
                className="py-3.5 pl-4 text-base text-white bg-purple-700/80"
                value={toLocation}
                onChangeText={setToLocation}
              />
            </BlurView>
          </View>

          <View className="mb-2">
            <Text className="text-base font-semibold text-gray-900 mb-4">Filtres</Text>
            <View className="flex-row justify-between">
              {transportModes.map((mode, i) => (
                <TouchableOpacity key={i} className="items-center">
                  <View className={`w-12 h-12 rounded-xl justify-center items-center mb-2 ${
                    mode.active ? 'bg-purple-700' : 'bg-gray-100'
                  } shadow-lg`}>
                    <Ionicons
                      name={mode.icon as any}
                      size={20}
                      color={mode.active ? '#fff' : '#9ca3af'}
                    />
                  </View>
                  <Text className={`text-xs font-medium ${
                    mode.active ? 'text-gray-900' : 'text-gray-400'
                  }`}>
                    {mode.label}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>
        </BlurView>
      )}
    </View>
  );
}
