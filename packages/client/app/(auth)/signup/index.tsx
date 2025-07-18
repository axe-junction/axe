import { Redirect } from 'expo-router';

export default function SignUpScreen() {
  // Redirect to main auth screen
  return <Redirect href="/(auth)" />;
}
