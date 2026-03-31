import React from "react";
import { ActivityIndicator, StyleSheet, View } from "react-native";
import { colors } from "@/theme";

interface LoadingSpinnerProps {
  size?: "small" | "large";
}

export function LoadingSpinner({ size = "large" }: LoadingSpinnerProps) {
  return (
    <View style={styles.center}>
      <ActivityIndicator testID="loading-spinner" size={size} color={colors.primary} />
    </View>
  );
}

const styles = StyleSheet.create({
  center: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: colors.background,
  },
});
