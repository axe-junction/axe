import { useState } from 'react';
import { Alert } from 'react-native';
import { router } from 'expo-router';

export interface AuthUser {
  id: string;
  email: string;
  fullName: string;
}

export const useAuth = () => {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const login = async (email: string, password: string) => {
    try {
      setIsLoading(true);
      
      // TODO: Replace with actual API call
      // For now, simulating login
      if (email && password) {
        const mockUser: AuthUser = {
          id: '1',
          email,
          fullName: 'User Name'
        };
        
        setUser(mockUser);
        Alert.alert('Succès', 'Connexion réussie!');
        
        // Navigate to main app
        router.replace('/map');
      }
    } catch (error) {
      Alert.alert('Erreur', 'Erreur lors de la connexion. Veuillez réessayer.');
      console.error('Login error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const signUp = async (fullName: string, email: string, password: string) => {
    try {
      setIsLoading(true);
      
      // TODO: Replace with actual API call
      // For now, simulating signup
      if (fullName && email && password) {
        const mockUser: AuthUser = {
          id: '1',
          email,
          fullName
        };
        
        setUser(mockUser);
        Alert.alert('Succès', 'Compte créé avec succès!');
        
        // Navigate to main app
        router.replace('/map');
      }
    } catch (error) {
      Alert.alert('Erreur', 'Erreur lors de la création du compte. Veuillez réessayer.');
      console.error('SignUp error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const googleAuth = async () => {
    try {
      setIsLoading(true);
      
      // TODO: Implement Google authentication
      Alert.alert('Info', 'Authentification Google sera implémentée prochainement.');
    } catch (error) {
      Alert.alert('Erreur', 'Erreur lors de l\'authentification Google.');
      console.error('Google auth error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const forgotPassword = async () => {
    try {
      // TODO: Implement forgot password functionality
      Alert.alert('Info', 'Fonctionnalité de récupération de mot de passe sera implémentée prochainement.');
    } catch (error) {
      Alert.alert('Erreur', 'Erreur lors de la demande de récupération.');
      console.error('Forgot password error:', error);
    }
  };

  const logout = () => {
    setUser(null);
    router.replace('/(auth)');
  };

  return {
    user,
    isLoading,
    login,
    signUp,
    googleAuth,
    forgotPassword,
    logout
  };
};
