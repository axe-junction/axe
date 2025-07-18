import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, Alert, ScrollView } from 'react-native';
import { Feather, AntDesign } from '@expo/vector-icons';

interface SignUpFormProps {
  onSignUp: (fullName: string, email: string, password: string) => void;
  onGoogleSignUp?: () => void;
  onLoginPress: () => void;
  isLoading?: boolean;
}

export default function SignUpForm({ 
  onSignUp, 
  onGoogleSignUp = () => {}, 
  onLoginPress, 
  isLoading = false 
}: SignUpFormProps) {
  const [email, setEmail] = useState('');
  const [fullName, setFullName] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const validateForm = () => {
    if (!fullName.trim()) {
      Alert.alert('Erreur', 'Veuillez saisir votre nom complet');
      return false;
    }
    if (!email.trim()) {
      Alert.alert('Erreur', 'Veuillez saisir votre adresse e-mail');
      return false;
    }
    if (!password) {
      Alert.alert('Erreur', 'Veuillez saisir un mot de passe');
      return false;
    }
    if (password !== confirmPassword) {
      Alert.alert('Erreur', 'Les mots de passe ne correspondent pas');
      return false;
    }
    if (password.length < 6) {
      Alert.alert('Erreur', 'Le mot de passe doit contenir au moins 6 caractères');
      return false;
    }
    return true;
  };

  const handleSignUp = () => {
    if (validateForm()) {
      onSignUp(fullName.trim(), email.trim(), password);
    }
  };

  const handleGoogleSignUp = () => {
    onGoogleSignUp();
  };

  return (
    <View className="flex-1 bg-white">
      <ScrollView 
        className="flex-1"
        contentContainerStyle={{ padding: 24, paddingTop: 64 }}
        showsVerticalScrollIndicator={false}
        keyboardShouldPersistTaps="handled"
      >
        <Text className="text-2xl font-bold text-center mb-1">Créer un compte Mapi</Text>
        <Text className="text-sm text-gray-500 text-center mb-8">
          Inscrivez-vous pour commencer à utiliser notre application.
        </Text>

      <Text className="text-sm font-medium mb-1">Nom complet *</Text>
      <TextInput
        className="border border-gray-300 rounded-xl px-4 py-3 mb-4"
        placeholder="Jean Dupont"
        onChangeText={setFullName}
        value={fullName}
        autoCapitalize="words"
        autoCorrect={false}
      />

      <Text className="text-sm font-medium mb-1">Adresse e-mail *</Text>
      <TextInput
        className="border border-gray-300 rounded-xl px-4 py-3 mb-4"
        placeholder="email@gmail.com"
        keyboardType="email-address"
        onChangeText={setEmail}
        value={email}
        autoCapitalize="none"
        autoCorrect={false}
      />

      <Text className="text-sm font-medium mb-1">Mot de passe *</Text>
      <View className="flex-row items-center border border-gray-300 rounded-xl px-4 mb-4">
        <TextInput
          className="flex-1 py-3"
          placeholder="******"
          secureTextEntry={!showPassword}
          onChangeText={setPassword}
          value={password}
          autoCapitalize="none"
          autoCorrect={false}
        />
        <TouchableOpacity onPress={() => setShowPassword(!showPassword)}>
          <Feather name={showPassword ? 'eye' : 'eye-off'} size={20} color="#888" />
        </TouchableOpacity>
      </View>

      <Text className="text-sm font-medium mb-1">Confirmer le mot de passe *</Text>
      <View className="flex-row items-center border border-gray-300 rounded-xl px-4 mb-6">
        <TextInput
          className="flex-1 py-3"
          placeholder="******"
          secureTextEntry={!showConfirmPassword}
          onChangeText={setConfirmPassword}
          value={confirmPassword}
          autoCapitalize="none"
          autoCorrect={false}
        />
        <TouchableOpacity onPress={() => setShowConfirmPassword(!showConfirmPassword)}>
          <Feather name={showConfirmPassword ? 'eye' : 'eye-off'} size={20} color="#888" />
        </TouchableOpacity>
      </View>

      <TouchableOpacity 
        className={`py-4 rounded-2xl mb-6 ${isLoading ? 'bg-gray-400' : 'bg-primary'}`}
        onPress={handleSignUp}
        disabled={isLoading}
      >
        <Text className="text-white font-bold text-center">
          {isLoading ? 'Inscription...' : "S'inscrire"}
        </Text>
      </TouchableOpacity>

      <Text className="text-center text-sm text-gray-600">Vous avez déjà un compte ?</Text>
      <TouchableOpacity onPress={onLoginPress}>
        <Text className="text-center text-primary text-sm font-medium mb-6">Se connecter ici</Text>
      </TouchableOpacity>

      <View className="flex-row items-center mb-6">
        <View className="flex-1 h-px bg-gray-300" />
        <Text className="mx-3 text-gray-500 text-sm">Ou continuez avec</Text>
        <View className="flex-1 h-px bg-gray-300" />
      </View>

      <TouchableOpacity 
        className="flex-row items-center justify-center border border-gray-300 py-3 rounded-xl"
        onPress={handleGoogleSignUp}
      >
        <AntDesign name="google" size={20} color="black" />
        <Text className="ml-2 text-sm">Google</Text>
      </TouchableOpacity>
      </ScrollView>
    </View>
  );
}
