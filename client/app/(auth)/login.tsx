import { Redirect } from 'expo-router';

export default function LoginScreen() {
  // Redirect to main auth screen
  return <Redirect href="/(auth)" />;
}
