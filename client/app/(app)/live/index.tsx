import React from 'react';
import { View, Text, StyleSheet, SafeAreaView, TouchableOpacity, ScrollView } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { BlurView } from 'expo-blur';

export default function LiveScreen() {
  const liveUpdates = [
    {
      id: 1,
      type: 'train',
      line: 'Ligne 1',
      status: 'En cours',
      delay: '2 min de retard',
      time: 'Il y a 1 min',
      color: '#ef4444'
    },
    {
      id: 2,
      type: 'bus',
      line: 'Tramway A',
      status: 'À l\'heure',
      delay: null,
      time: 'Il y a 3 min',
      color: '#10b981'
    },
    {
      id: 3,
      type: 'subway',
      line: 'Métro M1',
      status: 'Perturbation',
      delay: '5 min de retard',
      time: 'Il y a 5 min',
      color: '#f59e0b'
    },
    {
      id: 4,
      type: 'train',
      line: 'Ligne 2',
      status: 'À l\'heure',
      delay: null,
      time: 'Il y a 8 min',
      color: '#10b981'
    }
  ];

  const getTransportIcon = (type: string) => {
    switch (type) {
      case 'train':
        return 'train-outline';
      case 'bus':
        return 'bus-outline';
      case 'subway':
        return 'subway-outline';
      default:
        return 'location-outline';
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Informations en direct</Text>
        <TouchableOpacity style={styles.refreshButton}>
          <Ionicons name="refresh" size={24} color="#6b46c1" />
        </TouchableOpacity>
      </View>

      {/* Live Status Indicator */}
      <BlurView intensity={10} style={styles.statusCard}>
        <View style={styles.statusIndicator}>
          <View style={styles.liveDot} />
          <Text style={styles.statusText}>Mises à jour en temps réel</Text>
        </View>
      </BlurView>

      {/* Updates List */}
      <ScrollView style={styles.scrollContainer} showsVerticalScrollIndicator={false}>
        <Text style={styles.sectionTitle}>Dernières mises à jour</Text>
        
        {liveUpdates.map((update) => (
          <BlurView key={update.id} intensity={10} style={styles.updateCard}>
            <View style={styles.updateHeader}>
              <View style={styles.transportInfo}>
                <View style={[styles.iconContainer, { backgroundColor: update.color }]}>
                  <Ionicons 
                    name={getTransportIcon(update.type) as any} 
                    size={20} 
                    color="#fff" 
                  />
                </View>
                <View style={styles.lineInfo}>
                  <Text style={styles.lineName}>{update.line}</Text>
                  <Text style={styles.updateTime}>{update.time}</Text>
                </View>
              </View>
              <View style={[styles.statusBadge, { backgroundColor: update.color }]}>
                <Text style={styles.statusBadgeText}>{update.status}</Text>
              </View>
            </View>
            {update.delay && (
              <Text style={styles.delayText}>{update.delay}</Text>
            )}
          </BlurView>
        ))}

        {/* Emergency Alert */}
        <BlurView intensity={10} style={[styles.updateCard, styles.alertCard]}>
          <View style={styles.alertHeader}>
            <Ionicons name="warning" size={24} color="#ef4444" />
            <Text style={styles.alertTitle}>Alerte Réseau</Text>
          </View>
          <Text style={styles.alertMessage}>
            Travaux sur la ligne principale. Prévoir 10-15 minutes supplémentaires pour les trajets vers le centre-ville.
          </Text>
          <Text style={styles.alertTime}>Il y a 15 min</Text>
        </BlurView>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8fafc',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 20,
    paddingVertical: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#e5e7eb',
  },
  headerTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#1a202c',
  },
  refreshButton: {
    padding: 8,
  },
  statusCard: {
    marginHorizontal: 20,
    marginVertical: 16,
    padding: 16,
    borderRadius: 16,
    backgroundColor: 'rgba(255, 255, 255, 0.8)',
    overflow: 'hidden',
  },
  statusIndicator: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  liveDot: {
    width: 12,
    height: 12,
    borderRadius: 6,
    backgroundColor: '#ef4444',
    marginRight: 12,
    shadowColor: '#ef4444',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.8,
    shadowRadius: 4,
  },
  statusText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#1a202c',
  },
  scrollContainer: {
    flex: 1,
    paddingHorizontal: 20,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#1a202c',
    marginBottom: 16,
  },
  updateCard: {
    marginBottom: 12,
    padding: 16,
    borderRadius: 16,
    backgroundColor: 'rgba(255, 255, 255, 0.8)',
    overflow: 'hidden',
  },
  updateHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  transportInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  iconContainer: {
    width: 40,
    height: 40,
    borderRadius: 20,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  lineInfo: {
    flex: 1,
  },
  lineName: {
    fontSize: 16,
    fontWeight: '600',
    color: '#1a202c',
  },
  updateTime: {
    fontSize: 12,
    color: '#6b7280',
    marginTop: 2,
  },
  statusBadge: {
    paddingHorizontal: 12,
    paddingVertical: 4,
    borderRadius: 12,
  },
  statusBadgeText: {
    fontSize: 12,
    fontWeight: '600',
    color: '#fff',
  },
  delayText: {
    fontSize: 14,
    color: '#ef4444',
    fontWeight: '500',
  },
  alertCard: {
    borderLeftWidth: 4,
    borderLeftColor: '#ef4444',
  },
  alertHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 8,
  },
  alertTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#ef4444',
    marginLeft: 8,
  },
  alertMessage: {
    fontSize: 14,
    color: '#374151',
    lineHeight: 20,
    marginBottom: 8,
  },
  alertTime: {
    fontSize: 12,
    color: '#6b7280',
  },
});
