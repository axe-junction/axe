import React, { useState } from 'react';
import { Alert } from 'react-native';
import LoginForm from './LoginForm';
import SignUpForm from './SignUpForm';
import ForgotPasswordForm from './ForgotPasswordForm';

type AuthScreen = 'login' | 'signup' | 'forgot-password';

interface AuthManagerProps {
  onAuthSuccess?: (user: any) => void;
}

export default function AuthManager({ onAuthSuccess }: AuthManagerProps) {
  const [currentScreen, setCurrentScreen] = useState<AuthScreen>('login');
  const [isLoading, setIsLoading] = useState(false);

  const handleLogin = async (email: string, password: string) => {
    try {
      setIsLoading(true);
      
      // Check for specific credentials
      if (email.trim() !== 'soyed@gmail.com' || password !== 'soyed123') {
        Alert.alert('Erreur', 'Email ou mot de passe incorrect');
        return;
      }
      
      // Simulate login delay
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const mockUser = {
        id: '1',
        email: 'soyed@gmail.com',
        fullName: 'Soyed'
      };
      
      Alert.alert('Succès', 'Connexion réussie!');
      onAuthSuccess?.(mockUser);
    } catch (error) {
      Alert.alert('Erreur', 'Erreur lors de la connexion. Veuillez réessayer.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSignUp = async (fullName: string, email: string, password: string) => {
    try {
      setIsLoading(true);
      
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const mockUser = {
        id: '1',
        email,
        fullName
      };
      
      Alert.alert('Succès', 'Compte créé avec succès!');
      onAuthSuccess?.(mockUser);
    } catch (error) {
      Alert.alert('Erreur', 'Erreur lors de la création du compte. Veuillez réessayer.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleGoogleAuth = async () => {
    try {
      setIsLoading(true);
      
      Alert.alert('Info', 'Authentification Google sera implémentée prochainement.');
    } catch (error) {
      Alert.alert('Erreur', 'Erreur lors de l\'authentification Google.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleForgotPassword = () => {
    setCurrentScreen('forgot-password');
  };

  const handleResetPassword = async (email: string) => {
    try {
      setIsLoading(true);
      
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      Alert.alert(
        'E-mail envoyé',
        'Un lien de réinitialisation a été envoyé à votre adresse e-mail.',
        [
          {
            text: 'OK',
            onPress: () => setCurrentScreen('login')
          }
        ]
      );
    } catch (error) {
      Alert.alert('Erreur', 'Une erreur est survenue. Veuillez réessayer.');
    } finally {
      setIsLoading(false);
    }
  };

  const navigateToSignUp = () => {
    setCurrentScreen('signup');
  };

  const navigateToLogin = () => {
    setCurrentScreen('login');
  };

  const navigateBack = () => {
    setCurrentScreen('login');
  };

  switch (currentScreen) {
    case 'login':
      return (
        <LoginForm
          onLogin={handleLogin}
          onGoogleLogin={handleGoogleAuth}
          onForgotPassword={handleForgotPassword}
          onSignUpPress={navigateToSignUp}
          isLoading={isLoading}
        />
      );
    
    case 'signup':
      return (
        <SignUpForm
          onSignUp={handleSignUp}
          onGoogleSignUp={handleGoogleAuth}
          onLoginPress={navigateToLogin}
          isLoading={isLoading}
        />
      );
    
    case 'forgot-password':
      return (
        <ForgotPasswordForm
          onResetPassword={handleResetPassword}
          onBackPress={navigateBack}
          isLoading={isLoading}
        />
      );
    
    default:
      return (
        <LoginForm
          onLogin={handleLogin}
          onGoogleLogin={handleGoogleAuth}
          onForgotPassword={handleForgotPassword}
          onSignUpPress={navigateToSignUp}
          isLoading={isLoading}
        />
      );
  }
}
