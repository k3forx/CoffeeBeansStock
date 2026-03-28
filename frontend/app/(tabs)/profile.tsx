import { View, Text, TouchableOpacity, StyleSheet, Alert } from "react-native";
import { useAuthStore } from "../../src/stores/auth";
import { colors, typography, spacing, radius, shadows } from "@/theme";

export default function ProfileScreen() {
  const { user, logout } = useAuthStore();

  const handleLogout = () => {
    Alert.alert("ログアウト", "ログアウトしますか？", [
      { text: "キャンセル", style: "cancel" },
      {
        text: "ログアウト",
        style: "destructive",
        onPress: () => logout(),
      },
    ]);
  };

  return (
    <View style={styles.container}>
      <View style={styles.profileCard}>
        <View style={styles.avatar}>
          <Text style={styles.avatarText}>{user?.name?.charAt(0) ?? "?"}</Text>
        </View>
        <Text style={styles.name}>{user?.name ?? "ゲスト"}</Text>
        <Text style={styles.sub}>ゲストアカウント</Text>
      </View>

      <Text style={styles.sectionTitle}>設定</Text>
      <View style={styles.settingsCard}>
        <View style={styles.settingRow}>
          <Text style={styles.settingLabel}>在庫アラート</Text>
          <Text style={styles.settingValue}>
            {user?.low_stock_threshold ?? 100}g以下
          </Text>
        </View>
        <View style={[styles.settingRow, { borderBottomWidth: 0 }]}>
          <Text style={styles.settingLabel}>通知</Text>
          <Text style={styles.settingValue}>
            {user?.notification_enabled ? "ON" : "OFF"}
          </Text>
        </View>
      </View>

      <TouchableOpacity onPress={handleLogout}>
        <Text style={styles.logoutText}>ログアウト</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background, padding: spacing["2xl"] },
  profileCard: {
    backgroundColor: colors.surface,
    borderRadius: radius.lg,
    padding: spacing["3xl"],
    alignItems: "center",
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: colors.border,
    marginBottom: spacing["2xl"],
    ...shadows.sm,
  },
  avatar: {
    width: 88,
    height: 88,
    borderRadius: 44,
    borderWidth: 1.5,
    borderColor: colors.accent,
    justifyContent: "center",
    alignItems: "center",
    marginBottom: spacing.lg,
  },
  avatarText: {
    ...typography.titleLarge,
    color: colors.textPrimary,
    fontSize: 32,
    fontWeight: "300",
  },
  name: {
    ...typography.titleLarge,
    color: colors.textPrimary,
    marginBottom: spacing.xs,
  },
  sub: {
    ...typography.bodySmall,
    color: colors.textSecondary,
  },
  sectionTitle: {
    ...typography.labelMedium,
    color: colors.textSecondary,
    marginBottom: spacing.md,
    paddingLeft: spacing.xs,
  },
  settingsCard: {
    backgroundColor: colors.surface,
    borderRadius: radius.md,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: colors.border,
    overflow: "hidden",
    ...shadows.sm,
  },
  settingRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    paddingVertical: spacing.lg,
    paddingHorizontal: spacing.xl,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.borderLight,
  },
  settingLabel: {
    ...typography.bodyMedium,
    color: colors.textPrimary,
  },
  settingValue: {
    ...typography.bodyMedium,
    color: colors.textSecondary,
  },
  logoutText: {
    ...typography.bodyMedium,
    color: colors.danger,
    textAlign: "center",
    marginTop: spacing["4xl"],
  },
});
