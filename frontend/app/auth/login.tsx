import { useState, useRef, useEffect } from "react";
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  Alert,
  Animated,
} from "react-native";
import { MaterialCommunityIcons } from "@expo/vector-icons";
import { authApi } from "../../src/api/auth";
import { useAuthStore } from "../../src/stores/auth";
import { ApiError } from "../../src/api/client";
import { colors, typography, spacing, radius, shadows } from "@/theme";

export default function LoginScreen() {
  const [loading, setLoading] = useState(false);
  const setAuth = useAuthStore((s) => s.setAuth);

  const fadeAnims = useRef([...Array(4)].map(() => new Animated.Value(0))).current;
  const slideAnims = useRef([...Array(4)].map(() => new Animated.Value(10))).current;

  useEffect(() => {
    const animations = fadeAnims.map((anim, i) =>
      Animated.parallel([
        Animated.timing(anim, { toValue: 1, duration: 400, delay: i * 150, useNativeDriver: true }),
        Animated.timing(slideAnims[i], { toValue: 0, duration: 400, delay: i * 150, useNativeDriver: true }),
      ])
    );
    Animated.parallel(animations).start();
  }, []);

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
        <Animated.View style={{ opacity: fadeAnims[0], transform: [{ translateY: slideAnims[0] }] }}>
          <View style={styles.iconCircle}>
            <MaterialCommunityIcons name="coffee" size={40} color={colors.primary} />
          </View>
        </Animated.View>

        <Animated.View style={{ opacity: fadeAnims[1], transform: [{ translateY: slideAnims[1] }] }}>
          <Text style={styles.titleEn}>Coffee Beans Stock</Text>
          <Text style={styles.titleJa}>珈琲豆の記録</Text>
        </Animated.View>

        <Animated.View style={{ opacity: fadeAnims[2], transform: [{ translateY: slideAnims[2] }] }}>
          <Text style={styles.subtitle}>コーヒー豆の在庫を簡単管理</Text>
        </Animated.View>

        <Animated.View style={{ opacity: fadeAnims[3], transform: [{ translateY: slideAnims[3] }] }}>
          <TouchableOpacity
            style={[styles.button, loading && styles.buttonDisabled]}
            onPress={handleStart}
            disabled={loading}
          >
            <Text style={styles.buttonText}>
              {loading ? "準備中..." : "はじめる"}
            </Text>
          </TouchableOpacity>
        </Animated.View>
      </View>

      <Text style={styles.footer}>— 珈琲のある暮らし —</Text>
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
