import { useEffect, useState } from "react";
import { Slot, useRouter, useSegments } from "expo-router";
import { View, ActivityIndicator } from "react-native";
import { useAuthStore } from "../src/stores/auth";
import { authApi } from "../src/api/auth";

export default function RootLayout() {
  const isLoading = useAuthStore((s) => s.isLoading);
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const accessToken = useAuthStore((s) => s.accessToken);
  const segments = useSegments();
  const router = useRouter();
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
    useAuthStore.getState().loadStoredTokens();
  }, []);

  useEffect(() => {
    if (!isLoading && isAuthenticated && accessToken && !useAuthStore.getState().user) {
      authApi
        .getMe()
        .then((user) => {
          useAuthStore.setState({ user });
        })
        .catch(() => {
          useAuthStore.getState().logout();
        });
    }
  }, [isLoading, isAuthenticated, accessToken]);

  useEffect(() => {
    if (!isMounted || isLoading) return;
    const inAuthGroup = segments[0] === "auth";
    if (!isAuthenticated && !inAuthGroup) {
      router.replace("/auth/login");
    } else if (isAuthenticated && inAuthGroup) {
      router.replace("/(tabs)");
    }
  }, [isMounted, isLoading, isAuthenticated, segments]);

  if (isLoading) {
    return (
      <View style={{ flex: 1, justifyContent: "center", alignItems: "center", backgroundColor: "#F5F0E8" }}>
        <ActivityIndicator size="large" color="#4A3F35" />
      </View>
    );
  }

  return <Slot />;
}
