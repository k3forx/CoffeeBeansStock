import { useCallback, useState } from "react";
import { View, Text, TextInput, TouchableOpacity, StyleSheet, Alert } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { useFocusEffect } from "expo-router";
import { useAuthStore } from "../../src/stores/auth";
import { usersApi } from "../../src/api/users";
import { beansApi } from "../../src/api/beans";
import { showApiError } from "../../src/utils/errorHandler";
import { colors, typography, spacing, radius, shadows } from "@/theme";
import { Button } from "@/components/Button";

export default function ProfileScreen() {
  const { user, setUser } = useAuthStore();
  const [editingGrams, setEditingGrams] = useState(false);
  const [gramsInput, setGramsInput] = useState("");
  const [saving, setSaving] = useState(false);
  const [beanCount, setBeanCount] = useState<number | null>(null);

  useFocusEffect(
    useCallback(() => {
      beansApi.list(1000, 0).then(r => setBeanCount(r.beans.length)).catch(() => {});
    }, [])
  );

  const handleStartEditGrams = () => {
    setGramsInput(String(user?.grams_per_cup ?? 15));
    setEditingGrams(true);
  };

  const handleSaveGrams = async () => {
    const value = parseInt(gramsInput, 10);
    if (isNaN(value) || value < 1 || value > 100) {
      Alert.alert("エラー", "1〜100の範囲で入力してください");
      return;
    }
    setSaving(true);
    try {
      const updated = await usersApi.updateMe({ grams_per_cup: value });
      setUser(updated);
      setEditingGrams(false);
    } catch (e) {
      showApiError(e, "更新に失敗しました");
    } finally {
      setSaving(false);
    }
  };

  return (
    <View style={styles.container}>
      <View style={styles.profileCard}>
        <View style={styles.avatar}>
          <Text style={styles.avatarText}>{user?.name?.charAt(0) ?? "?"}</Text>
        </View>
        <Text style={styles.name}>{user?.name ?? "ゲスト"}</Text>
        <Text style={styles.sub}>ゲストアカウント</Text>
        {beanCount !== null && (
          <View style={styles.statBadge}>
            <Text style={styles.statText}>
              登録した珈琲豆: <Text style={styles.statNumber}>{beanCount}</Text> 種類
            </Text>
          </View>
        )}
      </View>

      <Text style={styles.sectionTitle}>設定</Text>
      <View style={styles.settingsCard}>
        <TouchableOpacity style={styles.settingRow} activeOpacity={0.7}>
          <Text style={styles.settingLabel}>在庫アラート</Text>
          <View style={styles.settingRight}>
            <Text style={styles.settingValue}>
              {user?.low_stock_threshold ?? 100}g以下
            </Text>
            <Ionicons name="chevron-forward" size={18} color={colors.textTertiary} />
          </View>
        </TouchableOpacity>
        <TouchableOpacity style={styles.settingRow} activeOpacity={0.7}>
          <Text style={styles.settingLabel}>通知</Text>
          <View style={styles.settingRight}>
            <Text style={styles.settingValue}>
              {user?.notification_enabled ? "ON" : "OFF"}
            </Text>
            <Ionicons name="chevron-forward" size={18} color={colors.textTertiary} />
          </View>
        </TouchableOpacity>
        <TouchableOpacity style={[styles.settingRow, { borderBottomWidth: 0 }]} activeOpacity={0.7}>
          <Text style={styles.settingLabel}>1杯あたりのグラム数</Text>
          {editingGrams ? (
            <View style={styles.editGramsRow}>
              <TextInput
                style={styles.gramsInput}
                value={gramsInput}
                onChangeText={setGramsInput}
                keyboardType="number-pad"
                autoFocus
                selectTextOnFocus
              />
              <Text style={styles.gramsUnit}>g</Text>
              <Button
                size="small"
                title="保存"
                onPress={handleSaveGrams}
                loading={saving}
                loadingText="..."
              />
              <Button
                size="small"
                variant="secondary"
                title="取消"
                onPress={() => setEditingGrams(false)}
                disabled={saving}
                style={{ borderWidth: 0 }}
              />
            </View>
          ) : (
            <TouchableOpacity onPress={handleStartEditGrams}>
              <Text style={styles.settingValueTappable}>
                {user?.grams_per_cup ?? 15}g
              </Text>
            </TouchableOpacity>
          )}
        </TouchableOpacity>
      </View>

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
  statBadge: {
    marginTop: spacing.lg,
    backgroundColor: "rgba(200,184,160,0.15)",
    paddingVertical: 5,
    paddingHorizontal: 14,
    borderRadius: radius.full,
    alignSelf: "center",
  },
  statText: {
    ...typography.bodySmall,
    color: colors.textSecondary,
  },
  statNumber: {
    fontWeight: "600",
    color: colors.textPrimary,
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
    alignItems: "center",
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
  settingRight: {
    flexDirection: "row",
    alignItems: "center",
    gap: 6,
  },
  settingValueTappable: {
    ...typography.bodyMedium,
    color: colors.accentDark,
    textDecorationLine: "underline",
  },
  editGramsRow: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.xs,
  },
  gramsInput: {
    width: 48,
    height: 32,
    borderBottomWidth: 1,
    borderBottomColor: colors.primary,
    textAlign: "center",
    ...typography.bodyMedium,
    color: colors.textPrimary,
    padding: 0,
  },
  gramsUnit: {
    ...typography.bodyMedium,
    color: colors.textSecondary,
  },
});
