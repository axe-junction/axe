import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, Alert, ScrollView } from 'react-native';
import { Ionicons } from '@expo/vector-icons';

interface ForgotPasswordFormProps {
  onResetPassword: (email: string) => void;
  onBackPress: () => void;
  isLoading?: boolean;
}

export default function ForgotPasswordForm({ onResetPassword, onBackPress, isLoading = false }: ForgotPasswordFormProps) {
  const [email, setEmail] = useState('');

  const handleResetPassword = () => {
    if (!email.trim()) {
      Alert.alert('Erreur', 'Veuillez saisir votre adresse e-mail');
      return;
    }
    onResetPassword(email);
  };

  return (
    <View className="flex-1 bg-white">
      <ScrollView 
        className="flex-1"
        contentContainerStyle={{ padding: 24, paddingTop: 64, flexGrow: 1 }}
        showsVerticalScrollIndicator={false}
        keyboardShouldPersistTaps="handled"
      >
        <TouchableOpacity 
          onPress={onBackPress}
          className="mb-8"
        >
          <Ionicons name="arrow-back" size={24} color="black" />
        </TouchableOpacity>

        <Text className="text-2xl font-bold text-center mb-1">Mot de passe oublié ?</Text>
        <Text className="text-sm text-gray-500 text-center mb-8">
          Saisissez votre adresse e-mail et nous vous enverrons un lien pour réinitialiser votre mot de passe.
        </Text>

      <Text className="text-sm font-medium mb-1">Adresse e-mail *</Text>
      <TextInput
        className="border border-gray-300 rounded-xl px-4 py-3 mb-6"
        placeholder="email@gmail.com"
        keyboardType="email-address"
        onChangeText={setEmail}
        value={email}
        autoCapitalize="none"
        autoCorrect={false}
      />

      <TouchableOpacity 
        className={`py-4 rounded-2xl mb-6 ${isLoading ? 'bg-gray-400' : 'bg-primary'}`}
        onPress={handleResetPassword}
        disabled={isLoading}
      >
        <Text className="text-white font-bold text-center">
          {isLoading ? 'Envoi en cours...' : 'Envoyer le lien'}
        </Text>
      </TouchableOpacity>

      <TouchableOpacity onPress={onBackPress}>
        <Text className="text-center text-primary text-sm font-medium">
          Retour à la connexion
        </Text>
      </TouchableOpacity>
      </ScrollView>
    </View>
  );
}
