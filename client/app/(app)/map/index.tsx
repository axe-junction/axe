import React, { useState } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, SafeAreaView, TextInput } from 'react-native';
import { Ionicons, MaterialCommunityIcons, Feather } from '@expo/vector-icons';
import MapView, { Marker, Polyline } from 'react-native-maps';
import { BlurView } from 'expo-blur';

export default function TransitApp() {
    const [currentScreen, setCurrentScreen] = useState<'map' | 'destination'>('map');
    const [searchQuery, setSearchQuery] = useState('');
    const [fromLocation, setFromLocation] = useState('');
    const [toLocation, setToLocation] = useState('');

    const handleSearchClick = () => setCurrentScreen('destination');
    const handleBackToMap = () => setCurrentScreen('map');

  const recentDestinations = [
    { name: "1er Mai", location: "P2WH+RFW, Alger" },
    { name: "Sonatrach - Direction Générale", location: "P2WH+RFW, Hydra" },
  ];

  const routeStations = [
    { 
      id: 1, 
      name: "Station Départ", 
      coordinate: { latitude: 36.7538, longitude: 3.0588 },
      type: "start"
    },
    { 
      id: 2, 
      name: "Station République", 
      coordinate: { latitude: 36.7580, longitude: 3.0520 },
      type: "intermediate"
    },
    { 
      id: 3, 
      name: "Station Centre-Ville", 
      coordinate: { latitude: 36.7620, longitude: 3.0450 },
      type: "intermediate"
    },
    { 
      id: 4, 
      name: "Station Université", 
      coordinate: { latitude: 36.7660, longitude: 3.0380 },
      type: "intermediate"
    },
    { 
      id: 5, 
      name: "Station Arrivée", 
      coordinate: { latitude: 36.7700, longitude: 3.0310 },
      type: "end"
    },
  ];

  const routePath = [
    { latitude: 36.7538, longitude: 3.0588 },
    { latitude: 36.7545, longitude: 3.0575 },
    { latitude: 36.7555, longitude: 3.0560 },
    { latitude: 36.7565, longitude: 3.0545 },
    { latitude: 36.7575, longitude: 3.0530 },
    { latitude: 36.7580, longitude: 3.0520 },
    { latitude: 36.7590, longitude: 3.0500 },
    { latitude: 36.7600, longitude: 3.0480 },
    { latitude: 36.7610, longitude: 3.0465 },
    { latitude: 36.7620, longitude: 3.0450 },
    { latitude: 36.7635, longitude: 3.0430 },
    { latitude: 36.7645, longitude: 3.0410 },
    { latitude: 36.7655, longitude: 3.0395 },
    { latitude: 36.7660, longitude: 3.0380 },
    { latitude: 36.7670, longitude: 3.0360 },
    { latitude: 36.7680, longitude: 3.0340 },
    { latitude: 36.7690, longitude: 3.0325 },
    { latitude: 36.7700, longitude: 3.0310 },
  ];

    const transportModes = [
      { icon: 'train-outline', label: "Train", active: true, iconSet: Ionicons },
      { icon: 'walk-outline', label: "Marche", active: false, iconSet: Ionicons },
      { icon: 'bus-outline', label: "Tramway", active: true, iconSet: Ionicons },
      { icon: 'subway-outline', label: "Métro", active: false, iconSet: Ionicons },
      { icon: 'shield-checkmark-outline', label: "Mode sûr", active: true, iconSet: Ionicons },
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

    if (currentScreen === 'destination') {
      return (
        <View style={styles.container}>
          <MapView
            style={styles.mapBackground}
            customMapStyle={customMapStyle}
            initialRegion={{
              latitude: 36.7538,
              longitude: 3.0588,
              latitudeDelta: 0.05,
              longitudeDelta: 0.05,
            }}
          >
            <Marker
              coordinate={{ latitude: 36.7538, longitude: 3.0588 }}
              pinColor="#6b46c1"
            />
            
          </MapView>

          <SafeAreaView style={styles.statusBar}>
            <Text style={styles.time}>9:41</Text>
            <View style={styles.signalBars}>
              <View style={[styles.signalBar, styles.signalBarStrong]} />
              <View style={[styles.signalBar, styles.signalBarStrong]} />
              <View style={[styles.signalBar, styles.signalBarStrong]} />
              <View style={[styles.signalBar, styles.signalBarWeak]} />
            </View>
          </SafeAreaView>

          <BlurView intensity={10} style={styles.destinationCard}>
            <View style={styles.cardHeader}>
              <TouchableOpacity onPress={handleBackToMap} style={styles.backButton}>
                <Ionicons name="arrow-back" size={24} color="#6b46c1" />
              </TouchableOpacity>
              <Text style={styles.cardTitle}>Choisissez votre destination</Text>
            </View>

            <View style={styles.inputContainer}>
              <BlurView intensity={10} style={styles.inputWrapper}>
                <View style={styles.inputDot} />
                <TextInput
                  placeholder="Départ"
                  placeholderTextColor="rgba(107, 114, 128, 1)"
                  style={styles.inputField}
                  value={fromLocation}
                  onChangeText={setFromLocation}
                />
              </BlurView>
              <BlurView intensity={10} style={styles.inputWrapper}>
                <TextInput
                  placeholder="Destination"
                  placeholderTextColor="rgba(255, 255, 255, 1)"
                  style={[styles.inputField, styles.destinationInput]}
                  value={toLocation}
                  onChangeText={setToLocation}
                />
              </BlurView>
            </View>

            <View style={styles.filterSection}>
              <Text style={styles.sectionTitle}>Filtres</Text>
              <View style={styles.transportModes}>
                {transportModes.map((mode, index) => {
                  const IconComponent = mode.iconSet;
                  return (
                    <TouchableOpacity 
                      key={index} 
                      style={styles.modeContainer}
                      onPress={() => console.log(`Selected ${mode.label}`)}
                    >
                      <View style={[
                        styles.modeButton,
                        mode.active ? styles.activeModeButton : styles.inactiveModeButton
                      ]}>
                        <IconComponent
                          name={mode.icon as any}
                          size={20}
                          color={mode.active ? '#fff' : '#9ca3af'}
                        />
                      </View>
                      <Text style={[
                        styles.modeLabel,
                        mode.active ? styles.activeModeLabel : styles.inactiveModeLabel
                      ]}>
                        {mode.label}
                      </Text>
                    </TouchableOpacity>
                  );
                })}
              </View>
            </View>
          </BlurView>
        </View>
      );
    }

    return (
      <SafeAreaView style={styles.container}>

        <MapView
          style={styles.mapBackground}
          customMapStyle={customMapStyle}
          initialRegion={{
            latitude: 36.7538,
            longitude: 3.0588,
            latitudeDelta: 0.05,
            longitudeDelta: 0.05,
          }}
        >
          <Marker
            coordinate={{ latitude: 36.7538, longitude: 3.0588 }}
            pinColor="#6b46c1"
            title="Station Départ"
            description="Départ"
          />
          <Marker
            coordinate={{ latitude: 36.7700, longitude: 3.0310 }}
            pinColor="#38a169"
            title="Station Arrivée"
            description="Arrivée"
          />
          <Polyline
            coordinates={routePath}
            strokeColor="#6b46c1"
            strokeWidth={4}
            lineDashPattern={[0]}
          />
        </MapView>

        <SafeAreaView style={styles.statusBar}>
          <Text style={styles.time}>9:41</Text>
          <View style={styles.signalBars}>
            <View style={[styles.signalBar, styles.signalBarStrong]} />
            <View style={[styles.signalBar, styles.signalBarStrong]} />
            <View style={[styles.signalBar, styles.signalBarStrong]} />
            <View style={[styles.signalBar, styles.signalBarWeak]} />
          </View>
        </SafeAreaView>

        <BlurView intensity={10} style={styles.offlineNotification}>
          <Ionicons name="wifi" size={16} color="#ed8936" />
          <Text style={styles.offlineText}>vous êtes actuellement hors ligne</Text>
        </BlurView>

        <TouchableOpacity 
          style={styles.searchButton}
          onPress={handleSearchClick}
        >
          <Text style={styles.searchButtonText}>Rechercher une destination</Text>
        </TouchableOpacity>

        <BlurView intensity={10} style={styles.recentDestinations}>
          <Text style={styles.sectionTitle}>Destinations récentes</Text>
          {recentDestinations.map((destination, index) => (
            <TouchableOpacity 
              key={index} 
              style={styles.destinationItem}
              onPress={() => {
                setToLocation(destination.name);
                handleSearchClick();
              }}
            >
              <View style={styles.destinationIcon}>
                <Ionicons name="location" size={20} color="#fff" />
              </View>
              <View style={styles.destinationText}>
                <Text style={styles.destinationName}>{destination.name}</Text>
                <Text style={styles.destinationLocation}>{destination.location}</Text>
              </View>
            </TouchableOpacity>
          ))}
        </BlurView>

        {/* <BlurView intensity={20} style={styles.bottomNavigation}>
          <TouchableOpacity style={[styles.navButton, styles.activeNavButton]}>
            <Ionicons name="map" size={24} color="#6b46c1" />
            <Text style={[styles.navButtonText, styles.activeNavButtonText]}>Map</Text>
          </TouchableOpacity>
          
          <TouchableOpacity style={styles.navButton}>
            <Ionicons name="radio" size={24} color="#9ca3af" />
            <Text style={styles.navButtonText}>Live</Text>
          </TouchableOpacity>
          
          <TouchableOpacity style={styles.navButton}>
            <Ionicons name="person" size={24} color="#9ca3af" />
            <Text style={styles.navButtonText}>Profil</Text>
          </TouchableOpacity>
        </BlurView> */}
      </SafeAreaView>
    );
  }

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      position: 'relative',
      paddingBottom: 90,
    },
    mapBackground: {
      ...StyleSheet.absoluteFillObject,
    },
    statusBar: {
      position: 'absolute',
      top: 0,
      left: 0,
      right: 0,
      padding: 16,
      flexDirection: 'row',
      justifyContent: 'space-between',
      alignItems: 'center',
    },
    time: {
      fontSize: 16,
      fontWeight: '600',
      color: '#000',
    },
    signalBars: {
      flexDirection: 'row',
      alignItems: 'center',
    },
    signalBar: {
      width: 3,
      height: 6,
      borderRadius: 2,
      marginRight: 2,
    },
    signalBarStrong: {
      backgroundColor: '#000',
    },
    signalBarWeak: {
      backgroundColor: '#9ca3af',
    },
    destinationCard: {
      position: 'absolute',
      top: 80,
      left: 20,
      right: 20,
      borderRadius: 24,
      padding: 20,
      backgroundColor: 'rgba(255, 255, 255, 0.7)',
      shadowColor: '#000',
      shadowOffset: { width: 0, height: 4 },
      shadowOpacity: 0.1,
      shadowRadius: 10,
      overflow: 'hidden',
    },
    cardHeader: {
      flexDirection: 'row',
      alignItems: 'center',
      marginBottom: 20,
    },
    backButton: {
      marginRight: 10,
    },
    cardTitle: {
      fontSize: 20,
      fontWeight: 'bold',
      color: '#1a202c',
    },
    inputContainer: {
      marginBottom: 20,
    },
    inputWrapper: {
      position: 'relative',
      marginBottom: 10,
      borderRadius: 12,
      overflow: 'hidden',
    },
    inputDot: {
      position: 'absolute',
      left: 12,
      top: 18,
      width: 6,
      height: 6,
      borderRadius: 3,
      backgroundColor: 'rgba(156, 163, 175, 0.8)',
      zIndex: 1,
    },
    inputField: {
      paddingVertical: 14,
      paddingLeft: 30,
      fontSize: 16,
      color: '#1a202c',
      backgroundColor: 'rgba(243, 244, 246, 0.8)',
    },
    destinationInput: {
      backgroundColor: 'rgba(107, 70, 193, 0.8)',
      color: '#fff',
      paddingLeft: 16,
    },
    filterSection: {
      marginBottom: 10,
    },
    sectionTitle: {
      fontSize: 16,
      fontWeight: '600',
      color: '#1a202c',
      marginBottom: 16,
    },
    transportModes: {
      flexDirection: 'row',
      justifyContent: 'space-between',
    },
    modeContainer: {
      alignItems: 'center',
    },
    modeButton: {
      width: 48,
      height: 48,
      borderRadius: 12,
      justifyContent: 'center',
      alignItems: 'center',
      marginBottom: 8,
      shadowColor: '#000',
      shadowOffset: { width: 0, height: 2 },
      shadowOpacity: 0.1,
      shadowRadius: 4,
    },
    activeModeButton: {
      backgroundColor: '#6b46c1',
    },
    inactiveModeButton: {
      backgroundColor: '#f3f4f6',
    },
    modeLabel: {
      fontSize: 12,
      fontWeight: '500',
    },
    activeModeLabel: {
      color: '#1a202c',
    },
    inactiveModeLabel: {
      color: '#9ca3af',
    },
    offlineNotification: {
      position: 'absolute',
      top: 80,
      alignSelf: 'center',
      borderRadius: 20,
      paddingVertical: 8,
      paddingHorizontal: 16,
      backgroundColor: 'rgba(255, 255, 255, 0.7)',
      flexDirection: 'row',
      alignItems: 'center',
      shadowColor: '#000',
      shadowOffset: { width: 0, height: 2 },
      shadowOpacity: 0.1,
      shadowRadius: 4,
      overflow: 'hidden',
    },
    offlineText: {
      marginLeft: 8,
      fontSize: 14,
      color: '#4a5568',
      fontWeight: '500',
    },
    searchButton: {
      position: 'absolute',
      top: '50%',
      left: 20,
      right: 20,
      backgroundColor: '#6b46c1',
      borderRadius: 24,
      paddingVertical: 16,
      alignItems: 'center',
      shadowColor: '#000',
      shadowOffset: { width: 0, height: 4 },
      shadowOpacity: 0.2,
      shadowRadius: 10,
    },
    searchButtonText: {
      color: '#fff',
      fontSize: 16,
      fontWeight: '600',
    },
    recentDestinations: {
      position: 'absolute',
      bottom: 100,
      left: 20,
      right: 20,
      borderRadius: 24,
      padding: 20,
      backgroundColor: 'rgba(255, 255, 255, 0.7)',
      shadowColor: '#000',
      shadowOffset: { width: 0, height: 4 },
      shadowOpacity: 0.1,
      shadowRadius: 10,
      overflow: 'hidden',
    },
    destinationItem: {
      flexDirection: 'row',
      alignItems: 'center',
      paddingVertical: 12,
    },
    destinationIcon: {
      width: 40,
      height: 40,
      borderRadius: 20,
      backgroundColor: '#6b46c1',
      justifyContent: 'center',
      alignItems: 'center',
      marginRight: 12,
    },
    destinationText: {
      flex: 1,
    },
    destinationName: {
      fontSize: 14,
      fontWeight: '600',
      color: '#1a202c',
    },
    destinationLocation: {
      fontSize: 12,
      color: '#d1d5db',
    },
    
    bottomNavigation: {
      position: 'absolute',
      bottom: 30,
      left: 30,
      right: 30,
      height: 70,
      borderRadius: 25,
      flexDirection: 'row',
      alignItems: 'center',
      justifyContent: 'space-around',
      backgroundColor: 'rgba(17, 24, 39, 0.8)',
      borderWidth: 1,
      borderColor: 'rgba(75, 85, 99, 0.4)',
      paddingHorizontal: 20,
      shadowColor: '#000',
      shadowOffset: {
        width: 0,
        height: 4,
      },
      shadowOpacity: 0.3,
      shadowRadius: 8,
      elevation: 8,
    },
    
    navButton: {
      alignItems: 'center',
      justifyContent: 'center',
      paddingVertical: 8,
      paddingHorizontal: 16,
      borderRadius: 16,
      minWidth: 70,
    },
    
    activeNavButton: {
      backgroundColor: 'rgba(107, 70, 193, 0.2)',
    },
    
    navButtonText: {
      fontSize: 11,
      color: '#9ca3af',
      marginTop: 4,
      fontWeight: '500',
    },
    
    activeNavButtonText: {
      color: '#6b46c1',
    },
  });


