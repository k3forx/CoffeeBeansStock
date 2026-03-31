import { useCallback, useEffect } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  Alert,
  ActivityIndicator,
  ScrollView,
  KeyboardAvoidingView,
  Platform,
} from "react-native";
import { Feather } from "@expo/vector-icons";
import { useFocusEffect, useLocalSearchParams, useRouter } from "expo-router";
import { useAuthStore } from "@/stores/auth";
import { useBeanDetail } from "@/hooks/useBeanDetail";
import { useBeanForm } from "@/hooks/useBeanForm";
import { useUsageTracking } from "@/hooks/useUsageTracking";
import { colors, typography, spacing, radius, shadows, getStockColor, formStyles } from "@/theme";
import { Button } from "@/components/Button";
import { ROAST_LEVELS, ROAST_DETAILS, ROAST_LEVEL_LABELS, ROAST_DETAIL_LABELS } from "@/constants/roastLevels";
import { ChipSelector } from "@/components/ChipSelector";
import { FormInput } from "@/components/FormInput";

export default function BeanDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const user = useAuthStore((s) => s.user);
  const gramsPerCup = user?.grams_per_cup ?? 15;
  const usagePresets = [
    { grams: gramsPerCup, label: "1杯分" },
    { grams: gramsPerCup * 2, label: "2杯分" },
  ];

  const detail = useBeanDetail();
  const form = useBeanForm();
  const usage = useUsageTracking({ beanId: id!, onBeanUpdated: detail.setBean });

  useFocusEffect(
    useCallback(() => {
      detail.loadBean(id!);
      usage.loadUsages();
    }, [id])
  );

  useEffect(() => {
    if (detail.bean) {
      form.resetToBean(detail.bean);
    }
  }, [detail.bean, detail.editing]);

  const handleSave = async () => {
    const result = form.validate();
    if (!result.valid) return;

    await detail.updateBean(id!, form.toUpdateInput(result.stock));
  };

  const handleDelete = () => {
    Alert.alert("削除確認", `「${detail.bean?.name}」を削除しますか？`, [
      { text: "キャンセル", style: "cancel" },
      {
        text: "削除",
        style: "destructive",
        onPress: async () => {
          const success = await detail.deleteBean(id!);
          if (success) router.back();
        },
      },
    ]);
  };

  const handleDeleteUsage = (usageId: string) => {
    Alert.alert("削除確認", "この使用記録を削除しますか？", [
      { text: "キャンセル", style: "cancel" },
      {
        text: "削除",
        style: "destructive",
        onPress: () => usage.handleDeleteUsage(usageId),
      },
    ]);
  };

  if (detail.loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }

  if (!detail.bean) return null;

  const bean = detail.bean;
  const stockColor = getStockColor(bean.current_stock);

  const roastDisplayText = (() => {
    const levelLabel = ROAST_LEVEL_LABELS[bean.roast_level];
    if (bean.roast_detail) {
      return `${levelLabel}（${ROAST_DETAIL_LABELS[bean.roast_detail]}）`;
    }
    return levelLabel;
  })();

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === "ios" ? "padding" : "height"}
    >
      {detail.editing ? (
        <View style={styles.editBanner}>
          <Feather name="edit-2" size={14} color={colors.warning} />
          <Text style={styles.editBannerText}>編集中</Text>
        </View>
      ) : (
        <View style={styles.editBar}>
          <TouchableOpacity style={styles.editPill} onPress={detail.startEditing}>
            <Feather name="edit-2" size={14} color={colors.accentDark} />
            <Text style={styles.editPillText}>編集</Text>
          </TouchableOpacity>
        </View>
      )}
      <ScrollView contentContainerStyle={styles.scroll} keyboardShouldPersistTaps="handled">
        {detail.editing ? (
          <>
            <FormInput
              label="名前"
              required
              value={form.fields.name}
              onChangeText={(v) => form.setField("name", v)}
              error={form.errors.name}
            />

            <FormInput
              label="産地"
              value={form.fields.origin}
              onChangeText={(v) => form.setField("origin", v)}
            />

            <Text style={formStyles.label}>焙煎度 <Text style={formStyles.required}>*</Text></Text>
            <ChipSelector
              items={ROAST_LEVELS}
              selected={form.fields.roastLevel}
              onSelect={form.handleRoastLevelSelect}
            />

            {form.fields.roastLevel !== "" && (
              <>
                <Text style={formStyles.label}>焙煎度（詳細）</Text>
                <ChipSelector
                  items={ROAST_DETAILS[form.fields.roastLevel]}
                  selected={form.fields.roastDetail}
                  onSelect={(v) => form.setField("roastDetail", form.fields.roastDetail === v ? "" : v)}
                />
              </>
            )}

            <FormInput
              label="在庫数 (g)"
              required
              value={form.fields.currentStock}
              onChangeText={(v) => form.setField("currentStock", v)}
              error={form.errors.currentStock}
              keyboardType="numeric"
            />

            <Text style={formStyles.label}>メモ</Text>
            <TextInput
              style={[styles.input, styles.textArea]}
              value={form.fields.notes}
              onChangeText={(v) => form.setField("notes", v)}
              multiline
              numberOfLines={4}
              placeholderTextColor={colors.textTertiary}
              underlineColorAndroid="transparent"
            />

            <View style={styles.editActions}>
              <Button
                variant="secondary"
                title="キャンセル"
                onPress={() => {
                  detail.cancelEditing();
                  form.resetToBean(bean);
                }}
                style={styles.editActionButton}
              />
              <Button
                title="保存"
                onPress={handleSave}
                loading={detail.saving}
                loadingText="保存中..."
                style={styles.editActionButton}
              />
            </View>
          </>
        ) : (
          <>
            <View style={styles.hero}>
              <Text style={styles.heroName}>{bean.name}</Text>
              <View style={[styles.heroCircle, { borderColor: stockColor }]}>
                <Text style={[styles.heroStockNum, { color: stockColor }]}>
                  {bean.current_stock}
                </Text>
                <Text style={styles.heroStockUnit}>gram</Text>
              </View>
            </View>

            <View style={styles.infoCard}>
              {bean.origin && (
                <View style={styles.infoRow}>
                  <Text style={styles.infoLabel}>産地</Text>
                  <Text style={styles.infoValue}>{bean.origin}</Text>
                </View>
              )}
              <View style={styles.infoRow}>
                <Text style={styles.infoLabel}>焙煎度</Text>
                <Text style={styles.infoValue}>{roastDisplayText}</Text>
              </View>
              {bean.notes && (
                <View style={styles.infoRow}>
                  <Text style={styles.infoLabel}>メモ</Text>
                  <Text style={styles.infoValue}>{bean.notes}</Text>
                </View>
              )}
              <View style={[styles.infoRow, styles.infoRowLast]}>
                <Text style={styles.infoLabel}>登録日</Text>
                <Text style={styles.infoValue}>
                  {new Date(bean.created_at).toLocaleDateString("ja-JP")}
                </Text>
              </View>
            </View>

            <View style={styles.usageSection}>
              <Text style={styles.sectionTitle}>使用記録</Text>

              <View style={styles.presetRow}>
                {usagePresets.map((preset, index) => (
                  <TouchableOpacity
                    key={preset.grams}
                    style={[
                      styles.presetBtn,
                      index > 0 && styles.presetBtnOutline,
                      usage.quickButtonLoading !== null && formStyles.buttonDisabled,
                    ]}
                    onPress={() => usage.handleQuickUsage(preset.grams)}
                    disabled={usage.quickButtonLoading !== null}
                  >
                    <Text style={[styles.presetBtnText, index > 0 && styles.presetBtnTextOutline]}>
                      {usage.quickButtonLoading === preset.grams ? "..." : `${preset.grams}g`}
                    </Text>
                    <Text style={[styles.presetBtnSub, index > 0 && styles.presetBtnSubOutline]}>
                      {preset.label}
                    </Text>
                  </TouchableOpacity>
                ))}
              </View>

              <View style={styles.manualInputRow}>
                <TextInput
                  style={styles.manualInput}
                  value={usage.manualGrams}
                  onChangeText={usage.setManualGrams}
                  keyboardType="numeric"
                  placeholder="グラム数"
                  placeholderTextColor={colors.textTertiary}
                  underlineColorAndroid="transparent"
                />
                <Text style={styles.manualInputUnit}>g</Text>
                <TouchableOpacity
                  style={[styles.manualButton, (usage.manualLoading || !usage.manualGrams) && formStyles.buttonDisabled]}
                  onPress={usage.handleManualUsage}
                  disabled={usage.manualLoading || !usage.manualGrams}
                >
                  <Text style={styles.manualButtonText}>
                    {usage.manualLoading ? "記録中..." : "記録"}
                  </Text>
                </TouchableOpacity>
              </View>

              <Text style={styles.historyTitle}>最近の使用履歴</Text>
              {usage.usageLoading ? (
                <ActivityIndicator size="small" color={colors.primary} />
              ) : usage.usages.length === 0 ? (
                <Text style={styles.emptyText}>まだ使用記録がありません</Text>
              ) : (
                usage.usages.map((u) => (
                  <View key={u.id} style={styles.historyItem}>
                    <View style={styles.historyInfo}>
                      <Text style={styles.historyDate}>
                        {new Date(u.usage_date).toLocaleDateString("ja-JP")}
                      </Text>
                      <Text style={styles.historyGrams}>{u.quantity}g</Text>
                    </View>
                    <TouchableOpacity
                      onPress={() => handleDeleteUsage(u.id)}
                      disabled={usage.deletingUsageId === u.id}
                    >
                      <Text style={styles.historyDeleteText}>
                        {usage.deletingUsageId === u.id ? "..." : "削除"}
                      </Text>
                    </TouchableOpacity>
                  </View>
                ))
              )}
            </View>
          </>
        )}
      </ScrollView>
      {!detail.editing && (
        <View style={styles.bottomBar}>
          <TouchableOpacity style={styles.deleteBtn} onPress={handleDelete}>
            <Feather name="trash-2" size={16} color={colors.danger} />
            <Text style={styles.deleteBtnText}>この豆を削除する</Text>
          </TouchableOpacity>
        </View>
      )}
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  center: { flex: 1, justifyContent: "center", alignItems: "center", backgroundColor: colors.background },
  editBar: {
    alignItems: "flex-end",
    paddingHorizontal: spacing.xl,
    paddingTop: spacing.md,
  },
  editPill: {
    flexDirection: "row",
    alignItems: "center",
    gap: 4,
    paddingVertical: 6,
    paddingHorizontal: 14,
    borderWidth: 1,
    borderColor: colors.accent,
    borderRadius: radius.full,
  },
  editPillText: {
    ...typography.labelMedium,
    color: colors.accentDark,
    letterSpacing: 0.5,
  },
  editBanner: {
    backgroundColor: colors.warningLight,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.xl,
    flexDirection: "row",
    alignItems: "center",
    gap: 6,
  },
  editBannerText: {
    ...typography.labelMedium,
    color: colors.warning,
    letterSpacing: 0.5,
  },
  scroll: { padding: spacing["2xl"] },
  hero: {
    alignItems: "center",
    marginBottom: spacing["3xl"],
    paddingBottom: spacing["2xl"],
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.border,
  },
  heroName: {
    ...typography.titleLarge,
    color: colors.textPrimary,
    marginBottom: spacing.xl,
    textAlign: "center",
  },
  heroCircle: {
    width: 110,
    height: 110,
    borderRadius: 55,
    borderWidth: 1.5,
    justifyContent: "center",
    alignItems: "center",
  },
  heroStockNum: {
    ...typography.stockNumberLarge,
  },
  heroStockUnit: {
    ...typography.bodySmall,
    color: colors.textSecondary,
    letterSpacing: 2,
    marginTop: spacing.xs,
  },
  infoCard: {
    backgroundColor: colors.surface,
    borderRadius: radius.md,
    padding: spacing.xs,
    paddingHorizontal: spacing.xl,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: colors.border,
    ...shadows.sm,
  },
  infoRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    paddingVertical: spacing.lg,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.borderLight,
  },
  infoRowLast: {
    borderBottomWidth: 0,
  },
  infoLabel: {
    ...typography.bodySmall,
    color: colors.textSecondary,
    letterSpacing: 0.5,
  },
  infoValue: {
    ...typography.bodyMedium,
    color: colors.textPrimary,
    flex: 1,
    textAlign: "right",
  },
  input: {
    borderWidth: 0,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
    paddingVertical: spacing.md,
    paddingHorizontal: 0,
    ...typography.bodyLarge,
    color: colors.textPrimary,
    backgroundColor: "transparent",
  },
  textArea: { height: 100, textAlignVertical: "top" },
  editActions: { flexDirection: "row", gap: spacing.md, marginTop: spacing["3xl"] },
  editActionButton: { flex: 1 },
  usageSection: {
    marginTop: spacing["3xl"],
  },
  sectionTitle: {
    ...typography.labelLarge,
    color: colors.textPrimary,
    marginBottom: spacing.lg,
  },
  presetRow: {
    flexDirection: "row" as const,
    gap: spacing.sm,
    marginBottom: spacing.lg,
  },
  presetBtn: {
    flex: 1,
    backgroundColor: colors.primary,
    borderRadius: radius.md,
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.sm,
    alignItems: "center" as const,
    borderWidth: 1.5,
    borderColor: colors.primary,
  },
  presetBtnOutline: {
    backgroundColor: "transparent",
    borderColor: colors.border,
  },
  presetBtnText: {
    ...typography.labelLarge,
    color: colors.textInverse,
  },
  presetBtnTextOutline: {
    color: colors.textPrimary,
  },
  presetBtnSub: {
    ...typography.labelSmall,
    color: colors.textInverse,
    opacity: 0.7,
    marginTop: 2,
  },
  presetBtnSubOutline: {
    color: colors.textSecondary,
    opacity: 1,
  },
  manualInputRow: {
    flexDirection: "row" as const,
    alignItems: "center" as const,
    gap: spacing.sm,
    marginBottom: spacing["2xl"],
  },
  manualInput: {
    flex: 1,
    height: 48,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
    ...typography.bodyLarge,
    color: colors.textPrimary,
  },
  manualInputUnit: {
    ...typography.bodyMedium,
    color: colors.textSecondary,
  },
  manualButton: {
    backgroundColor: colors.primary,
    borderRadius: radius.sm,
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.xl,
  },
  manualButtonText: {
    ...typography.labelLarge,
    color: colors.textInverse,
  },
  historyTitle: {
    ...typography.labelMedium,
    color: colors.textSecondary,
    marginBottom: spacing.md,
  },
  emptyText: {
    ...typography.bodySmall,
    color: colors.textTertiary,
    textAlign: "center" as const,
    paddingVertical: spacing.xl,
  },
  historyItem: {
    flexDirection: "row" as const,
    justifyContent: "space-between" as const,
    alignItems: "center" as const,
    paddingVertical: spacing.md,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.borderLight,
  },
  historyInfo: {
    flexDirection: "row" as const,
    alignItems: "center" as const,
    gap: spacing.md,
  },
  historyDate: {
    ...typography.bodyMedium,
    color: colors.textPrimary,
  },
  historyGrams: {
    ...typography.labelLarge,
    color: colors.textPrimary,
  },
  historyDeleteText: {
    ...typography.bodySmall,
    color: colors.danger,
  },
  bottomBar: {
    padding: spacing.lg,
    alignItems: "center",
  },
  deleteBtn: {
    flexDirection: "row" as const,
    alignItems: "center" as const,
    justifyContent: "center" as const,
    gap: 6,
    paddingVertical: 10,
    paddingHorizontal: spacing.xl,
    borderWidth: 1,
    borderColor: colors.danger,
    borderRadius: radius.md,
  },
  deleteBtnText: {
    ...typography.bodySmall,
    color: colors.danger,
    letterSpacing: 0.5,
  },
});
