import React from 'react';
import { AuthManager } from '../../presentation/components/auth';
import { router } from 'expo-router';

export default function AuthIndex() {
  const handleAuthSuccess = (user: any) => {
    // Navigate to main app after successful authentication
    router.replace('/map');
  };

  return <AuthManager onAuthSuccess={handleAuthSuccess} />;
}
