import React from 'react';
import { View, ActivityIndicator, Text } from 'react-native';

interface LoadingSpinnerProps {
  message?: string;
  size?: 'small' | 'large';
}

export default function LoadingSpinner({ message = 'Chargement...', size = 'large' }: LoadingSpinnerProps) {
  return (
    <View className="flex-1 justify-center items-center bg-white">
      <ActivityIndicator size={size} color="#007AFF" />
      {message && (
        <Text className="mt-4 text-gray-600 text-center">{message}</Text>
      )}
    </View>
  );
}
