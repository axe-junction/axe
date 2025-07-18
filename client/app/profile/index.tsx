import React from 'react';
import { View, Text, StyleSheet, SafeAreaView, TouchableOpacity, ScrollView, Image } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { BlurView } from 'expo-blur';

export default function ProfileScreen() {
  const profileData = {
    name: "Ahmed Ben Ali",
    email: "ahmed.benali@email.com",
    phone: "+213 555 123 456",
    joinDate: "Membre depuis mars 2024",
    tripsCount: 47,
    favoriteRoutes: 5,
    points: 1250
  };

  const menuItems = [
    { icon: 'card-outline', title: 'Modes de paiement', subtitle: 'Gérer vos cartes' },
    { icon: 'time-outline', title: 'Historique des trajets', subtitle: `${profileData.tripsCount} trajets effectués` },
    { icon: 'heart-outline', title: 'Destinations favorites', subtitle: `${profileData.favoriteRoutes} destinations sauvegardées` },
    { icon: 'notifications-outline', title: 'Notifications', subtitle: 'Alertes et mises à jour' },
    { icon: 'shield-outline', title: 'Sécurité et confidentialité', subtitle: 'Paramètres de sécurité' },
    { icon: 'help-circle-outline', title: 'Aide et support', subtitle: 'FAQ et contact' },
    { icon: 'settings-outline', title: 'Paramètres', subtitle: 'Préférences de l\'application' },
  ];

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Profil</Text>
        <TouchableOpacity style={styles.editButton}>
          <Ionicons name="create-outline" size={24} color="#6b46c1" />
        </TouchableOpacity>
      </View>

      <ScrollView style={styles.scrollContainer} showsVerticalScrollIndicator={false}>
        {/* Profile Card */}
        <BlurView intensity={10} style={styles.profileCard}>
          <View style={styles.profileHeader}>
            <View style={styles.avatarContainer}>
              <Ionicons name="person" size={40} color="#fff" />
            </View>
            <View style={styles.profileInfo}>
              <Text style={styles.profileName}>{profileData.name}</Text>
              <Text style={styles.profileEmail}>{profileData.email}</Text>
              <Text style={styles.profileJoinDate}>{profileData.joinDate}</Text>
            </View>
          </View>
          
          {/* Stats */}
          <View style={styles.statsContainer}>
            <View style={styles.statItem}>
              <Text style={styles.statNumber}>{profileData.tripsCount}</Text>
              <Text style={styles.statLabel}>Trajets</Text>
            </View>
            <View style={styles.statDivider} />
            <View style={styles.statItem}>
              <Text style={styles.statNumber}>{profileData.favoriteRoutes}</Text>
              <Text style={styles.statLabel}>Favoris</Text>
            </View>
            <View style={styles.statDivider} />
            <View style={styles.statItem}>
              <Text style={styles.statNumber}>{profileData.points}</Text>
              <Text style={styles.statLabel}>Points</Text>
            </View>
          </View>
        </BlurView>

        {/* Quick Actions */}
        <View style={styles.quickActions}>
          <BlurView intensity={10} style={styles.actionCard}>
            <Ionicons name="qr-code" size={24} color="#6b46c1" />
            <Text style={styles.actionText}>QR Code</Text>
          </BlurView>
          <BlurView intensity={10} style={styles.actionCard}>
            <Ionicons name="wallet-outline" size={24} color="#6b46c1" />
            <Text style={styles.actionText}>Portefeuille</Text>
          </BlurView>
          <BlurView intensity={10} style={styles.actionCard}>
            <Ionicons name="ticket-outline" size={24} color="#6b46c1" />
            <Text style={styles.actionText}>Billets</Text>
          </BlurView>
        </View>

        {/* Menu Items */}
        <View style={styles.menuSection}>
          <Text style={styles.sectionTitle}>Paramètres du compte</Text>
          {menuItems.map((item, index) => (
            <TouchableOpacity key={index} style={styles.menuItem}>
              <BlurView intensity={10} style={styles.menuItemContent}>
                <View style={styles.menuItemLeft}>
                  <View style={styles.menuIcon}>
                    <Ionicons name={item.icon as any} size={20} color="#6b46c1" />
                  </View>
                  <View style={styles.menuText}>
                    <Text style={styles.menuTitle}>{item.title}</Text>
                    <Text style={styles.menuSubtitle}>{item.subtitle}</Text>
                  </View>
                </View>
                <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
              </BlurView>
            </TouchableOpacity>
          ))}
        </View>

        {/* Logout Button */}
        <TouchableOpacity style={styles.logoutButton}>
          <Text style={styles.logoutText}>Se déconnecter</Text>
        </TouchableOpacity>

        {/* App Version */}
        <Text style={styles.versionText}>Version 1.0.0</Text>
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
  editButton: {
    padding: 8,
  },
  scrollContainer: {
    flex: 1,
    paddingHorizontal: 20,
  },
  profileCard: {
    marginVertical: 20,
    padding: 20,
    borderRadius: 16,
    backgroundColor: 'rgba(255, 255, 255, 0.8)',
    overflow: 'hidden',
  },
  profileHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 20,
  },
  avatarContainer: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: '#6b46c1',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 16,
  },
  profileInfo: {
    flex: 1,
  },
  profileName: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#1a202c',
    marginBottom: 4,
  },
  profileEmail: {
    fontSize: 14,
    color: '#6b7280',
    marginBottom: 2,
  },
  profileJoinDate: {
    fontSize: 12,
    color: '#9ca3af',
  },
  statsContainer: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    alignItems: 'center',
    paddingTop: 20,
    borderTopWidth: 1,
    borderTopColor: '#e5e7eb',
  },
  statItem: {
    alignItems: 'center',
  },
  statNumber: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#6b46c1',
  },
  statLabel: {
    fontSize: 12,
    color: '#6b7280',
    marginTop: 4,
  },
  statDivider: {
    width: 1,
    height: 40,
    backgroundColor: '#e5e7eb',
  },
  quickActions: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 30,
  },
  actionCard: {
    flex: 1,
    alignItems: 'center',
    padding: 16,
    marginHorizontal: 4,
    borderRadius: 12,
    backgroundColor: 'rgba(255, 255, 255, 0.8)',
    overflow: 'hidden',
  },
  actionText: {
    fontSize: 12,
    color: '#1a202c',
    marginTop: 8,
    fontWeight: '500',
  },
  menuSection: {
    marginBottom: 30,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#1a202c',
    marginBottom: 16,
  },
  menuItem: {
    marginBottom: 8,
  },
  menuItemContent: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: 16,
    borderRadius: 12,
    backgroundColor: 'rgba(255, 255, 255, 0.8)',
    overflow: 'hidden',
  },
  menuItemLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  menuIcon: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: 'rgba(107, 70, 193, 0.1)',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  menuText: {
    flex: 1,
  },
  menuTitle: {
    fontSize: 16,
    fontWeight: '500',
    color: '#1a202c',
    marginBottom: 2,
  },
  menuSubtitle: {
    fontSize: 12,
    color: '#6b7280',
  },
  logoutButton: {
    backgroundColor: '#ef4444',
    padding: 16,
    borderRadius: 12,
    alignItems: 'center',
    marginBottom: 20,
  },
  logoutText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  versionText: {
    textAlign: 'center',
    fontSize: 12,
    color: '#9ca3af',
    marginBottom: 20,
  },
});
