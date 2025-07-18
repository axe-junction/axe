import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, Alert, ScrollView } from 'react-native';
import { Feather, AntDesign } from '@expo/vector-icons';

interface LoginFormProps {
  onLogin: (email: string, password: string) => void;
  onGoogleLogin?: () => void;
  onForgotPassword: () => void;
  onSignUpPress: () => void;
  isLoading?: boolean;
}

export default function LoginForm({ 
  onLogin, 
  onGoogleLogin = () => {}, 
  onForgotPassword, 
  onSignUpPress, 
  isLoading = false 
}: LoginFormProps) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  const handleLogin = () => {
    if (!email.trim()) {
      Alert.alert('Erreur', 'Veuillez saisir votre adresse e-mail');
      return;
    }
    if (!password) {
      Alert.alert('Erreur', 'Veuillez saisir votre mot de passe');
      return;
    }
    onLogin(email.trim(), password);
  };

  const handleGoogleLogin = () => {
    onGoogleLogin();
  };

  return (
    <View className="flex-1 bg-white">
      <ScrollView 
        className="flex-1"
        contentContainerStyle={{ padding: 24, paddingTop: 64, flexGrow: 1, justifyContent: 'center' }}
        showsVerticalScrollIndicator={false}
        keyboardShouldPersistTaps="handled"
      >
        <Text className="text-2xl font-bold text-center mb-1">Se connecter à Mapi</Text>
        <Text className="text-sm text-gray-500 text-center mb-8">
          Créez un compte ou connectez-vous pour découvrir notre application.
        </Text>

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

      <View className="items-end mb-6">
        <TouchableOpacity onPress={onForgotPassword}>
          <Text className="text-primary text-sm">Mot de passe oublié ?</Text>
        </TouchableOpacity>
      </View>

      <TouchableOpacity 
        className={`py-4 rounded-2xl mb-6 ${isLoading ? 'bg-gray-400' : 'bg-primary'}`}
        onPress={handleLogin}
        disabled={isLoading}
      >
        <Text className="text-white font-bold text-center">
          {isLoading ? 'Connexion...' : 'Se connecter'}
        </Text>
      </TouchableOpacity>

      <Text className="text-center text-sm text-gray-600">Vous n'avez pas de compte ?</Text>
      <TouchableOpacity onPress={onSignUpPress}>
        <Text className="text-center text-primary text-sm font-medium mb-6">S'inscrire ici</Text>
      </TouchableOpacity>

      <View className="flex-row items-center mb-6">
        <View className="flex-1 h-px bg-gray-300" />
        <Text className="mx-3 text-gray-500 text-sm">Ou continuez avec</Text>
        <View className="flex-1 h-px bg-gray-300" />
      </View>

      <TouchableOpacity 
        className="flex-row items-center justify-center border border-gray-300 py-3 rounded-xl"
        onPress={handleGoogleLogin}
      >
        <AntDesign name="google" size={20} color="black" />
        <Text className="ml-2 text-sm">Google</Text>
      </TouchableOpacity>
      </ScrollView>
    </View>
  );
}
