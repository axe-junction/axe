import React from 'react';
import { View, Text, StyleSheet, SafeAreaView, TouchableOpacity, ScrollView } from 'react-native';
import { BlurView } from 'expo-blur';
import { Ionicons } from '@expo/vector-icons';
import { router } from 'expo-router';

export default function LiveScreen() {
  const navigateToHome = () => {
    router.push('/map');
  };

  const navigateToProfile = () => {
    router.push('/profile');
  };

  const vtcData = {
    "yassir": {
      "price": 871,
      "eta": "15 min",
      "icon": "car-sport"
    },
    "InDrive": {
      "priceRange": "725–1015", 
      "eta": "15 min",
      "icon": "car"
    },
    "yango": {
      "price": 842,
      "eta": "15 min", 
      "icon": "car-outline"
    },
    "heetch": {
      "price": 900,
      "eta": "15 min",
      "icon": "car-sport-outline"
    }
  };

  const liveUpdates = [
    {
      id: 1,
      type: 'traffic',
      message: 'Embouteillage signalé sur Rue Hassiba Ben Bouali',
      time: 'Il y a 2 min',
      severity: 'high'
    },
    {
      id: 2,
      type: 'accident',
      message: 'Accident mineur à la Place des Martyrs',
      time: 'Il y a 5 min',
      severity: 'medium'
    },
    {
      id: 3,
      type: 'construction',
      message: 'Travaux en cours Boulevard Mohamed V',
      time: 'Il y a 10 min',
      severity: 'low'
    },
    {
      id: 4,
      type: 'weather',
      message: 'Conditions météo favorables à Alger',
      time: 'Il y a 15 min',
      severity: 'low'
    }
  ];

  const getIcon = (type: string) => {
    switch (type) {
      case 'traffic': return 'car';
      case 'accident': return 'warning';
      case 'construction': return 'construct';
      case 'weather': return 'partly-sunny';
      default: return 'information-circle';
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'high': return '#ef4444';
      case 'medium': return '#f59e0b';
      case 'low': return '#10b981';
      default: return '#6b7280';
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>En direct</Text>
        <Text style={styles.subtitle}>VTC Services & Live Updates</Text>
      </View>

      <ScrollView style={styles.content} showsVerticalScrollIndicator={false}>
        {/* VTC Services Cards */}
        <View style={styles.content}>
          <Text style={{fontSize: 18, fontWeight: 'bold', marginBottom: 16, color: '#1a202c'}}>Services VTC Disponibles</Text>
          {Object.entries(vtcData).slice(0, 2).map(([serviceName, data]) => (
            <BlurView key={serviceName} intensity={15} style={styles.updateCard}>
              <View style={styles.updateHeader}>
                <View style={styles.iconContainer}>
                  <Ionicons name={data.icon as any} size={20} color="#6b46c1" />
                </View>
                <Text style={styles.updateTime}>{data.eta}</Text>
              </View>
              <Text style={styles.updateMessage}>
                {serviceName}: {(data as any).price ? `${(data as any).price} DA` : (data as any).priceRange + ' DA'}
              </Text>
            </BlurView>
          ))}
        </View>

        {/* Live Updates */}
        {liveUpdates.map((update) => (
          <BlurView key={update.id} intensity={15} style={styles.updateCard}>
            <View style={styles.updateHeader}>
              <View style={[styles.iconContainer, { backgroundColor: getSeverityColor(update.severity) + '20' }]}>
                <Ionicons 
                  name={getIcon(update.type) as any} 
                  size={20} 
                  color={getSeverityColor(update.severity)} 
                />
              </View>
              <Text style={styles.updateTime}>{update.time}</Text>
            </View>
            <Text style={styles.updateMessage}>{update.message}</Text>
          </BlurView>
        ))}

        <TouchableOpacity style={styles.refreshButton}>
          <BlurView intensity={20} style={styles.refreshButtonInner}>
            <Ionicons name="refresh" size={20} color="#6b46c1" />
            <Text style={styles.refreshText}>Actualiser</Text>
          </BlurView>
        </TouchableOpacity>
      </ScrollView>

      <BlurView intensity={20} style={styles.bottomNavigation}>
        <TouchableOpacity style={styles.navButton} onPress={navigateToHome}>
          <Ionicons name="home" size={24} color="#9ca3af" />
          <Text style={styles.navButtonText}>Home</Text>
        </TouchableOpacity>
        
        <TouchableOpacity style={[styles.navButton, styles.activeNavButton]}>
          <Ionicons name="radio" size={24} color="#6b46c1" />
          <Text style={[styles.navButtonText, styles.activeNavButtonText]}>Live</Text>
        </TouchableOpacity>
        
        <TouchableOpacity style={styles.navButton} onPress={navigateToProfile}>
          <Ionicons name="person" size={24} color="#9ca3af" />
          <Text style={styles.navButtonText}>Profile</Text>
        </TouchableOpacity>
      </BlurView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8fafc',
  },
  
  header: {
    paddingTop: 20,
    paddingBottom: 20,
    paddingHorizontal: 20,
    backgroundColor: '#ffffff',
    borderBottomWidth: 1,
    borderBottomColor: '#e5e7eb',
  },
  
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#1a202c',
    marginBottom: 4,
  },
  
  subtitle: {
    fontSize: 14,
    color: '#6b7280',
  },
  
  content: {
    flex: 1,
    padding: 20,
    paddingBottom: 100,
  },
  
  updateCard: {
    backgroundColor: 'rgba(255, 255, 255, 0.8)',
    borderRadius: 16,
    padding: 16,
    marginBottom: 12,
    borderWidth: 1,
    borderColor: 'rgba(107, 114, 128, 0.1)',
    overflow: 'hidden',
  },
  
  updateHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: 8,
  },
  
  iconContainer: {
    width: 36,
    height: 36,
    borderRadius: 18,
    alignItems: 'center',
    justifyContent: 'center',
  },
  
  updateTime: {
    fontSize: 12,
    color: '#6b7280',
  },
  
  updateMessage: {
    fontSize: 14,
    color: '#1a202c',
    lineHeight: 20,
  },
  
  refreshButton: {
    marginTop: 20,
    marginBottom: 20,
  },
  
  refreshButtonInner: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 16,
    paddingHorizontal: 24,
    backgroundColor: 'rgba(255, 255, 255, 0.8)',
    borderRadius: 16,
    borderWidth: 1,
    borderColor: 'rgba(107, 70, 193, 0.2)',
    overflow: 'hidden',
  },
  
  refreshText: {
    fontSize: 16,
    color: '#6b46c1',
    fontWeight: '600',
    marginLeft: 8,
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
    overflow: 'hidden',
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