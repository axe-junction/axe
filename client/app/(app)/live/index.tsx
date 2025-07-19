import React, { useState } from 'react';
import { View, Text, StyleSheet, SafeAreaView, TouchableOpacity, ScrollView, Image } from 'react-native';
import { BlurView } from 'expo-blur';
import { Ionicons } from '@expo/vector-icons';
import { router } from 'expo-router';

interface VTCService {
  price?: number;
  priceRange?: string;
  eta: string;
  icon: string;
}

export default function LiveScreen() {
  const [currentView, setCurrentView] = useState<'cards' | 'list'>('cards');

  const navigateToHome = () => {
    router.push('/map');
  };

  const navigateToProfile = () => {
    router.push('/profile');
  };

  const vtcData: Record<string, VTCService> = {
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

  const renderVTCCards = () => (
    <View style={styles.cardsContainer}>
      <Text style={styles.sectionTitle}>Services VTC Disponibles</Text>
      <View style={styles.cardsGrid}>
        {Object.entries(vtcData).slice(0, 2).map(([serviceName, data], index) => (
          <BlurView key={serviceName} intensity={15} style={styles.vtcCard}>
            <View style={styles.cardHeader}>
              <View style={styles.carIcon}>
                <Ionicons name={data.icon as any} size={24} color="#6b46c1" />
              </View>
              <Text style={styles.serviceName}>{serviceName}</Text>
            </View>
            <View style={styles.cardBody}>
              <Text style={styles.priceText}>
                {(data as any).price ? `${(data as any).price} DA` : (data as any).priceRange + ' DA'}
              </Text>
              <View style={styles.etaContainer}>
                <Ionicons name="time-outline" size={16} color="#6b7280" />
                <Text style={styles.etaText}>{data.eta}</Text>
              </View>
            </View>
            <TouchableOpacity style={styles.bookButton}>
              <Text style={styles.bookButtonText}>Réserver</Text>
            </TouchableOpacity>
          </BlurView>
        ))}
      </View>
    </View>
  );

  const renderVTCList = () => (
    <View style={styles.listContainer}>
      <Text style={styles.sectionTitle}>Tous les Services VTC</Text>
      {Object.entries(vtcData).map(([serviceName, data], index) => (
        <BlurView key={serviceName} intensity={15} style={styles.vtcListItem}>
          <View style={styles.listItemLeft}>
            <View style={styles.listCarIcon}>
              <Ionicons name={data.icon as any} size={20} color="#6b46c1" />
            </View>
            <View style={styles.serviceInfo}>
              <Text style={styles.listServiceName}>{serviceName}</Text>
              <View style={styles.listEtaContainer}>
                <Ionicons name="time-outline" size={14} color="#6b7280" />
                <Text style={styles.listEtaText}>{data.eta}</Text>
              </View>
            </View>
          </View>
          <View style={styles.listItemRight}>
            <Text style={styles.listPriceText}>
              {(data as any).price ? `${(data as any).price} DA` : (data as any).priceRange + ' DA'}
            </Text>
            <TouchableOpacity style={styles.listBookButton}>
              <Text style={styles.listBookButtonText}>Réserver</Text>
            </TouchableOpacity>
          </View>
        </BlurView>
      ))}
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>En direct</Text>
        <Text style={styles.subtitle}>Services et mises à jour</Text>
        
        <View style={styles.toggleContainer}>
          <TouchableOpacity 
            style={[styles.toggleButton, currentView === 'cards' && styles.activeToggle]}
            onPress={() => setCurrentView('cards')}
          >
            <Ionicons name="grid-outline" size={20} color={currentView === 'cards' ? '#fff' : '#6b7280'} />
            <Text style={[styles.toggleText, currentView === 'cards' && styles.activeToggleText]}>Cards</Text>
          </TouchableOpacity>
          <TouchableOpacity 
            style={[styles.toggleButton, currentView === 'list' && styles.activeToggle]}
            onPress={() => setCurrentView('list')}
          >
            <Ionicons name="list-outline" size={20} color={currentView === 'list' ? '#fff' : '#6b7280'} />
            <Text style={[styles.toggleText, currentView === 'list' && styles.activeToggleText]}>List</Text>
          </TouchableOpacity>
        </View>
      </View>

      <ScrollView style={styles.content} showsVerticalScrollIndicator={false}>
        {currentView === 'cards' ? renderVTCCards() : renderVTCList()}
        
        <View style={styles.updatesSection}>
          <Text style={styles.sectionTitle}>Mises à jour en temps réel</Text>
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
        </View>

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
    marginBottom: 16,
  },

  toggleContainer: {
    flexDirection: 'row',
    backgroundColor: '#f3f4f6',
    borderRadius: 12,
    padding: 4,
  },

  toggleButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 8,
    paddingHorizontal: 16,
    borderRadius: 8,
  },

  activeToggle: {
    backgroundColor: '#6b46c1',
  },

  toggleText: {
    fontSize: 14,
    fontWeight: '500',
    color: '#6b7280',
    marginLeft: 6,
  },

  activeToggleText: {
    color: '#fff',
  },
  
  content: {
    flex: 1,
    padding: 20,
    paddingBottom: 100,
  },

  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#1a202c',
    marginBottom: 16,
  },

  // Cards View Styles
  cardsContainer: {
    marginBottom: 30,
  },

  cardsGrid: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    flexWrap: 'wrap',
  },

  vtcCard: {
    width: '48%',
    backgroundColor: 'rgba(255, 255, 255, 0.9)',
    borderRadius: 16,
    padding: 16,
    marginBottom: 12,
    borderWidth: 1,
    borderColor: 'rgba(107, 114, 128, 0.1)',
    overflow: 'hidden',
  },

  cardHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 12,
  },

  carIcon: {
    width: 32,
    height: 32,
    marginRight: 8,
    borderRadius: 8,
    backgroundColor: 'rgba(107, 70, 193, 0.1)',
    alignItems: 'center',
    justifyContent: 'center',
  },

  serviceName: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#1a202c',
    textTransform: 'capitalize',
  },

  cardBody: {
    marginBottom: 16,
  },

  priceText: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#6b46c1',
    marginBottom: 8,
  },

  etaContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },

  etaText: {
    fontSize: 14,
    color: '#6b7280',
    marginLeft: 4,
  },

  bookButton: {
    backgroundColor: '#6b46c1',
    borderRadius: 8,
    paddingVertical: 10,
    alignItems: 'center',
  },

  bookButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },

  // List View Styles
  listContainer: {
    marginBottom: 30,
  },

  vtcListItem: {
    backgroundColor: 'rgba(255, 255, 255, 0.9)',
    borderRadius: 16,
    padding: 16,
    marginBottom: 12,
    borderWidth: 1,
    borderColor: 'rgba(107, 114, 128, 0.1)',
    overflow: 'hidden',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },

  listItemLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },

  listCarIcon: {
    width: 40,
    height: 40,
    marginRight: 12,
    borderRadius: 8,
    backgroundColor: 'rgba(107, 70, 193, 0.1)',
    alignItems: 'center',
    justifyContent: 'center',
  },

  serviceInfo: {
    flex: 1,
  },

  listServiceName: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#1a202c',
    textTransform: 'capitalize',
    marginBottom: 4,
  },

  listEtaContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },

  listEtaText: {
    fontSize: 12,
    color: '#6b7280',
    marginLeft: 4,
  },

  listItemRight: {
    alignItems: 'flex-end',
  },

  listPriceText: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#6b46c1',
    marginBottom: 8,
  },

  listBookButton: {
    backgroundColor: '#6b46c1',
    borderRadius: 6,
    paddingVertical: 6,
    paddingHorizontal: 12,
  },

  listBookButtonText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: '600',
  },

  // Updates Section
  updatesSection: {
    marginBottom: 20,
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