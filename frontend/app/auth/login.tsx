import { useState } from "react";
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  Alert,
} from "react-native";
import { authApi } from "../../src/api/auth";
import { useAuthStore } from "../../src/stores/auth";
import { ApiError } from "../../src/api/client";
import { colors, typography, spacing, radius, shadows } from "@/theme";

export default function LoginScreen() {
  const [loading, setLoading] = useState(false);
  const setAuth = useAuthStore((s) => s.setAuth);

  const handleStart = async () => {
    setLoading(true);
    try {
      const result = await authApi.registerAnonymous();
      await setAuth(result.user, result.access_token, result.refresh_token);
    } catch (e) {
      const msg = e instanceof ApiError ? e.message : String(e);
      Alert.alert("エラー", msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <View style={styles.container}>
      <View style={styles.content}>
        <View style={styles.iconCircle}>
          <Text style={styles.iconEmoji}>☕</Text>
        </View>

        <Text style={styles.titleEn}>Coffee Beans Stock</Text>
        <Text style={styles.titleJa}>珈琲豆の記録</Text>
        <Text style={styles.subtitle}>コーヒー豆の在庫を簡単管理</Text>

        <TouchableOpacity
          style={[styles.button, loading && styles.buttonDisabled]}
          onPress={handleStart}
          disabled={loading}
        >
          <Text style={styles.buttonText}>
            {loading ? "準備中..." : "はじめる"}
          </Text>
        </TouchableOpacity>
      </View>

      <Text style={styles.footer}>— craft your cup —</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  content: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    padding: spacing["2xl"],
    paddingBottom: 80,
  },
  iconCircle: {
    width: 100,
    height: 100,
    borderRadius: 50,
    borderWidth: 1.5,
    borderColor: colors.accent,
    justifyContent: "center",
    alignItems: "center",
    marginBottom: spacing["3xl"],
  },
  iconEmoji: { fontSize: 36 },
  titleEn: {
    ...typography.displaySub,
    color: colors.textPrimary,
    marginBottom: spacing.sm,
  },
  titleJa: {
    ...typography.displayLarge,
    color: colors.textPrimary,
    marginBottom: spacing.sm,
  },
  subtitle: {
    ...typography.bodyMedium,
    color: colors.textSecondary,
    marginBottom: spacing["5xl"],
  },
  button: {
    backgroundColor: colors.primary,
    borderRadius: radius.sm,
    paddingVertical: spacing.lg,
    paddingHorizontal: spacing["3xl"],
    width: "80%",
    alignItems: "center",
    ...shadows.md,
  },
  buttonDisabled: { opacity: 0.6 },
  buttonText: {
    ...typography.labelLarge,
    color: colors.textInverse,
    letterSpacing: 2,
  },
  footer: {
    ...typography.bodySmall,
    color: colors.textTertiary,
    textAlign: "center",
    paddingBottom: spacing["5xl"],
    fontStyle: "italic",
  },
});
