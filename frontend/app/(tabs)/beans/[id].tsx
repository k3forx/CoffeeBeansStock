import { useCallback, useState } from "react";
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
import { beansApi } from "../../../src/api/beans";
import { usagesApi } from "../../../src/api/usages";
import { useAuthStore } from "../../../src/stores/auth";
import { showApiError } from "../../../src/utils/errorHandler";
import type { CoffeeBean, RoastLevel, RoastDetail, UsageHistory } from "../../../src/types/api";
import { colors, typography, spacing, radius, shadows, getStockColor, formStyles } from "@/theme";
import { ROAST_LEVELS, ROAST_DETAILS, ROAST_LEVEL_LABELS, ROAST_DETAIL_LABELS } from "../../../src/constants/roastLevels";
import { ChipSelector } from "../../../src/components/ChipSelector";
import { FormInput } from "../../../src/components/FormInput";
import { validateBeanForm } from "../../../src/utils/validation";

export default function BeanDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const user = useAuthStore((s) => s.user);
  const gramsPerCup = user?.grams_per_cup ?? 15;
  const usagePresets = [
    { grams: gramsPerCup, label: "1杯分" },
    { grams: gramsPerCup * 2, label: "2杯分" },
  ];
  const [bean, setBean] = useState<CoffeeBean | null>(null);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [saving, setSaving] = useState(false);

  const [name, setName] = useState("");
  const [origin, setOrigin] = useState("");
  const [roastLevel, setRoastLevel] = useState<RoastLevel | "">("");
  const [roastDetail, setRoastDetail] = useState<RoastDetail | "">("");
  const [currentStock, setCurrentStock] = useState("");
  const [notes, setNotes] = useState("");

  const [usages, setUsages] = useState<UsageHistory[]>([]);
  const [usageLoading, setUsageLoading] = useState(false);
  const [quickButtonLoading, setQuickButtonLoading] = useState<number | null>(null);
  const [manualGrams, setManualGrams] = useState("");
  const [manualLoading, setManualLoading] = useState(false);
  const [deletingUsageId, setDeletingUsageId] = useState<string | null>(null);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const loadBean = useCallback(async () => {
    try {
      const data = await beansApi.get(id);
      setBean(data);
      setName(data.name);
      setOrigin(data.origin ?? "");
      setRoastLevel(data.roast_level);
      setRoastDetail(data.roast_detail ?? "");
      setCurrentStock(String(data.current_stock));
      setNotes(data.notes ?? "");
    } catch {
      Alert.alert("エラー", "データの取得に失敗しました", [{ text: "OK", onPress: () => router.back() }]);
    } finally {
      setLoading(false);
    }
  }, [id, router]);

  const loadUsages = useCallback(async () => {
    setUsageLoading(true);
    try {
      const data = await usagesApi.list(id, 10, 0);
      setUsages(data.usages);
    } catch {
      // Silent fail — usage list is supplementary
    } finally {
      setUsageLoading(false);
    }
  }, [id]);

  useFocusEffect(
    useCallback(() => {
      loadBean();
      loadUsages();
    }, [loadBean, loadUsages])
  );

  const getTodayDate = () => {
    const d = new Date();
    const y = d.getFullYear();
    const m = String(d.getMonth() + 1).padStart(2, "0");
    const day = String(d.getDate()).padStart(2, "0");
    return `${y}-${m}-${day}`;
  };

  const refreshAfterUsage = async () => {
    const [updatedBean, updatedUsages] = await Promise.all([
      beansApi.get(id),
      usagesApi.list(id, 10, 0),
    ]);
    setBean(updatedBean);
    setCurrentStock(String(updatedBean.current_stock));
    setUsages(updatedUsages.usages);
  };

  const handleQuickUsage = async (grams: number) => {
    setQuickButtonLoading(grams);
    try {
      await usagesApi.create(id, {
        usage_date: getTodayDate(),
        quantity: grams,
      });
      await refreshAfterUsage();
    } catch (e) {
      showApiError(e, "記録に失敗しました", {
        INSUFFICIENT_STOCK: "在庫が不足しています",
      });
    } finally {
      setQuickButtonLoading(null);
    }
  };

  const handleManualUsage = async () => {
    const grams = parseInt(manualGrams, 10);
    if (isNaN(grams) || grams <= 0) {
      Alert.alert("エラー", "1g以上の数値を入力してください");
      return;
    }
    setManualLoading(true);
    try {
      await usagesApi.create(id, {
        usage_date: getTodayDate(),
        quantity: grams,
      });
      await refreshAfterUsage();
      setManualGrams("");
    } catch (e) {
      showApiError(e, "記録に失敗しました", {
        INSUFFICIENT_STOCK: "在庫が不足しています",
      });
    } finally {
      setManualLoading(false);
    }
  };

  const handleDeleteUsage = (usageId: string) => {
    Alert.alert("削除確認", "この使用記録を削除しますか？", [
      { text: "キャンセル", style: "cancel" },
      {
        text: "削除",
        style: "destructive",
        onPress: async () => {
          setDeletingUsageId(usageId);
          try {
            await usagesApi.delete(id, usageId);
            await refreshAfterUsage();
          } catch {
            Alert.alert("エラー", "削除に失敗しました");
          } finally {
            setDeletingUsageId(null);
          }
        },
      },
    ]);
  };

  const handleRoastLevelSelect = (level: RoastLevel) => {
    setRoastLevel(level);
    setRoastDetail("");
  };

  const handleSave = async () => {
    const result = validateBeanForm({ name, roastLevel, currentStock });
    if (!result.valid) {
      setErrors(result.errors);
      return;
    }
    setErrors({});

    setSaving(true);
    try {
      const updated = await beansApi.update(id, {
        name,
        origin: origin || undefined,
        roast_level: roastLevel as RoastLevel,
        roast_detail: roastDetail || undefined,
        current_stock: result.stock,
        notes: notes || undefined,
      });
      setBean(updated);
      setEditing(false);
    } catch (e) {
      showApiError(e, "更新に失敗しました");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = () => {
    Alert.alert("削除確認", `「${bean?.name}」を削除しますか？`, [
      { text: "キャンセル", style: "cancel" },
      {
        text: "削除",
        style: "destructive",
        onPress: async () => {
          try {
            await beansApi.delete(id);
            router.back();
          } catch {
            Alert.alert("エラー", "削除に失敗しました");
          }
        },
      },
    ]);
  };

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }

  if (!bean) return null;

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
      {editing ? (
        <View style={styles.editBanner}>
          <Feather name="edit-2" size={14} color={colors.warning} />
          <Text style={styles.editBannerText}>編集中</Text>
        </View>
      ) : (
        <View style={styles.editBar}>
          <TouchableOpacity style={styles.editPill} onPress={() => setEditing(true)}>
            <Feather name="edit-2" size={14} color={colors.accentDark} />
            <Text style={styles.editPillText}>編集</Text>
          </TouchableOpacity>
        </View>
      )}
      <ScrollView contentContainerStyle={styles.scroll} keyboardShouldPersistTaps="handled">
        {editing ? (
          <>
            <FormInput
              label="名前"
              required
              value={name}
              onChangeText={(v) => { setName(v); setErrors(prev => { const { name: _, ...rest } = prev; return rest; }); }}
              error={errors.name}
            />

            <FormInput
              label="産地"
              value={origin}
              onChangeText={setOrigin}
            />

            <Text style={formStyles.label}>焙煎度 <Text style={formStyles.required}>*</Text></Text>
            <ChipSelector
              items={ROAST_LEVELS}
              selected={roastLevel}
              onSelect={handleRoastLevelSelect}
            />

            {roastLevel !== "" && (
              <>
                <Text style={formStyles.label}>焙煎度（詳細）</Text>
                <ChipSelector
                  items={ROAST_DETAILS[roastLevel]}
                  selected={roastDetail}
                  onSelect={(v) => setRoastDetail(roastDetail === v ? "" : v)}
                />
              </>
            )}

            <FormInput
              label="在庫数 (g)"
              required
              value={currentStock}
              onChangeText={(v) => { setCurrentStock(v); setErrors(prev => { const { currentStock: _, ...rest } = prev; return rest; }); }}
              error={errors.currentStock}
              keyboardType="numeric"
            />

            <Text style={formStyles.label}>メモ</Text>
            <TextInput
              style={[styles.input, styles.textArea]}
              value={notes}
              onChangeText={setNotes}
              multiline
              numberOfLines={4}
              placeholderTextColor={colors.textTertiary}
              underlineColorAndroid="transparent"
            />

            <View style={styles.editActions}>
              <TouchableOpacity
                style={[styles.button, styles.cancelButton]}
                onPress={() => {
                  setEditing(false);
                  setName(bean.name);
                  setOrigin(bean.origin ?? "");
                  setRoastLevel(bean.roast_level);
                  setRoastDetail(bean.roast_detail ?? "");
                  setCurrentStock(String(bean.current_stock));
                  setNotes(bean.notes ?? "");
                }}
              >
                <Text style={styles.cancelButtonText}>キャンセル</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={[styles.button, styles.saveButton, saving && formStyles.buttonDisabled]}
                onPress={handleSave}
                disabled={saving}
              >
                <Text style={styles.saveButtonText}>{saving ? "保存中..." : "保存"}</Text>
              </TouchableOpacity>
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
                      quickButtonLoading !== null && formStyles.buttonDisabled,
                    ]}
                    onPress={() => handleQuickUsage(preset.grams)}
                    disabled={quickButtonLoading !== null}
                  >
                    <Text style={[styles.presetBtnText, index > 0 && styles.presetBtnTextOutline]}>
                      {quickButtonLoading === preset.grams ? "..." : `${preset.grams}g`}
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
                  value={manualGrams}
                  onChangeText={setManualGrams}
                  keyboardType="numeric"
                  placeholder="グラム数"
                  placeholderTextColor={colors.textTertiary}
                  underlineColorAndroid="transparent"
                />
                <Text style={styles.manualInputUnit}>g</Text>
                <TouchableOpacity
                  style={[styles.manualButton, (manualLoading || !manualGrams) && formStyles.buttonDisabled]}
                  onPress={handleManualUsage}
                  disabled={manualLoading || !manualGrams}
                >
                  <Text style={styles.manualButtonText}>
                    {manualLoading ? "記録中..." : "記録"}
                  </Text>
                </TouchableOpacity>
              </View>

              <Text style={styles.historyTitle}>最近の使用履歴</Text>
              {usageLoading ? (
                <ActivityIndicator size="small" color={colors.primary} />
              ) : usages.length === 0 ? (
                <Text style={styles.emptyText}>まだ使用記録がありません</Text>
              ) : (
                usages.map((usage) => (
                  <View key={usage.id} style={styles.historyItem}>
                    <View style={styles.historyInfo}>
                      <Text style={styles.historyDate}>
                        {new Date(usage.usage_date).toLocaleDateString("ja-JP")}
                      </Text>
                      <Text style={styles.historyGrams}>{usage.quantity}g</Text>
                    </View>
                    <TouchableOpacity
                      onPress={() => handleDeleteUsage(usage.id)}
                      disabled={deletingUsageId === usage.id}
                    >
                      <Text style={styles.historyDeleteText}>
                        {deletingUsageId === usage.id ? "..." : "削除"}
                      </Text>
                    </TouchableOpacity>
                  </View>
                ))
              )}
            </View>
          </>
        )}
      </ScrollView>
      {!editing && (
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

  // View mode - Hero
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

  // View mode - Info card
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

  // Edit mode
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
  button: { ...formStyles.buttonBase, flex: 1 },
  cancelButton: {
    backgroundColor: "transparent",
    borderWidth: 1,
    borderColor: colors.border,
  },
  cancelButtonText: {
    ...typography.labelLarge,
    color: colors.textSecondary,
  },
  saveButton: { backgroundColor: colors.primary },
  saveButtonText: {
    ...formStyles.buttonTextBase,
  },

  // Usage section
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

  // Bottom
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
