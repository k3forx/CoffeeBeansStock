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
import { useFocusEffect, useLocalSearchParams, useRouter } from "expo-router";
import { beansApi } from "../../../src/api/beans";
import { ApiError } from "../../../src/api/client";
import type { CoffeeBean } from "../../../src/types/api";
import { colors, typography, spacing, radius, shadows, getStockColor } from "@/theme";

export default function BeanDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const [bean, setBean] = useState<CoffeeBean | null>(null);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [saving, setSaving] = useState(false);

  const [name, setName] = useState("");
  const [origin, setOrigin] = useState("");
  const [roastLevel, setRoastLevel] = useState("");
  const [currentStock, setCurrentStock] = useState("");
  const [notes, setNotes] = useState("");

  const loadBean = async () => {
    try {
      const data = await beansApi.get(id);
      setBean(data);
      setName(data.name);
      setOrigin(data.origin ?? "");
      setRoastLevel(data.roast_level ?? "");
      setCurrentStock(String(data.current_stock));
      setNotes(data.notes ?? "");
    } catch {
      Alert.alert("エラー", "データの取得に失敗しました", [{ text: "OK", onPress: () => router.back() }]);
    } finally {
      setLoading(false);
    }
  };

  useFocusEffect(
    useCallback(() => {
      loadBean();
    }, [id])
  );

  const handleSave = async () => {
    if (!name) {
      Alert.alert("エラー", "名前は必須です");
      return;
    }
    const stock = parseInt(currentStock, 10);
    if (isNaN(stock) || stock < 0) {
      Alert.alert("エラー", "在庫数は0以上の数値を入力してください");
      return;
    }

    setSaving(true);
    try {
      const updated = await beansApi.update(id, {
        name,
        origin: origin || undefined,
        roast_level: roastLevel || undefined,
        current_stock: stock,
        notes: notes || undefined,
      });
      setBean(updated);
      setEditing(false);
    } catch (e) {
      const msg = e instanceof ApiError ? e.message : "更新に失敗しました";
      Alert.alert("エラー", msg);
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

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === "ios" ? "padding" : "height"}
    >
      {!editing && (
        <View style={styles.editBar}>
          <TouchableOpacity onPress={() => setEditing(true)}>
            <Text style={styles.editLink}>編集</Text>
          </TouchableOpacity>
        </View>
      )}
      <ScrollView contentContainerStyle={styles.scroll} keyboardShouldPersistTaps="handled">
        {editing ? (
          <>
            <Text style={styles.label}>名前 *</Text>
            <TextInput
              style={styles.input}
              value={name}
              onChangeText={setName}
              placeholderTextColor={colors.textTertiary}
              underlineColorAndroid="transparent"
            />

            <Text style={styles.label}>産地</Text>
            <TextInput
              style={styles.input}
              value={origin}
              onChangeText={setOrigin}
              placeholderTextColor={colors.textTertiary}
              underlineColorAndroid="transparent"
            />

            <Text style={styles.label}>焙煎度</Text>
            <TextInput
              style={styles.input}
              value={roastLevel}
              onChangeText={setRoastLevel}
              placeholderTextColor={colors.textTertiary}
              underlineColorAndroid="transparent"
            />

            <Text style={styles.label}>在庫数 (g) *</Text>
            <TextInput
              style={styles.input}
              value={currentStock}
              onChangeText={setCurrentStock}
              keyboardType="numeric"
              placeholderTextColor={colors.textTertiary}
              underlineColorAndroid="transparent"
            />

            <Text style={styles.label}>メモ</Text>
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
                  setRoastLevel(bean.roast_level ?? "");
                  setCurrentStock(String(bean.current_stock));
                  setNotes(bean.notes ?? "");
                }}
              >
                <Text style={styles.cancelButtonText}>キャンセル</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={[styles.button, styles.saveButton, saving && styles.buttonDisabled]}
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
              {bean.roast_level && (
                <View style={styles.infoRow}>
                  <Text style={styles.infoLabel}>焙煎度</Text>
                  <Text style={styles.infoValue}>{bean.roast_level}</Text>
                </View>
              )}
              {bean.notes && (
                <View style={styles.infoRow}>
                  <Text style={styles.infoLabel}>メモ</Text>
                  <Text style={styles.infoValue}>{bean.notes}</Text>
                </View>
              )}
              <View style={[styles.infoRow, { borderBottomWidth: 0 }]}>
                <Text style={styles.infoLabel}>登録日</Text>
                <Text style={styles.infoValue}>
                  {new Date(bean.created_at).toLocaleDateString("ja-JP")}
                </Text>
              </View>
            </View>
          </>
        )}
      </ScrollView>
      {!editing && (
        <View style={styles.bottomBar}>
          <TouchableOpacity onPress={handleDelete}>
            <Text style={styles.deleteText}>この豆を削除する</Text>
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
  editLink: {
    ...typography.labelLarge,
    color: colors.accentDark,
    letterSpacing: 1,
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
  label: {
    ...typography.labelMedium,
    color: colors.textSecondary,
    marginBottom: spacing.sm,
    marginTop: spacing.xl,
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
  button: { flex: 1, borderRadius: radius.sm, padding: spacing.lg, alignItems: "center" },
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
  buttonDisabled: { opacity: 0.6 },
  saveButtonText: {
    ...typography.labelLarge,
    color: colors.textInverse,
  },

  // Bottom
  bottomBar: {
    padding: spacing.lg,
    alignItems: "center",
  },
  deleteText: {
    ...typography.bodySmall,
    color: colors.danger,
    textDecorationLine: "underline",
    letterSpacing: 1,
  },
});
