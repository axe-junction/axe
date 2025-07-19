import { Stack } from "expo-router";
import "../presentation/styles/global.css";

export default function RootLayout() {
  return (
    <Stack
      screenOptions={{
        headerShown: false, 
      }}
    />
  );
}
