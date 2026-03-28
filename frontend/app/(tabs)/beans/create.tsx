import { useState } from "react";
import {
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  Alert,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
} from "react-native";
import { useRouter } from "expo-router";
import { beansApi } from "../../../src/api/beans";
import { ApiError } from "../../../src/api/client";
import type { RoastLevel } from "../../../src/types/api";
import { colors, typography, spacing, radius, shadows } from "@/theme";

export default function CreateBeanScreen() {
  const [name, setName] = useState("");
  const [origin, setOrigin] = useState("");
  const [roastLevel, setRoastLevel] = useState("");
  const [currentStock, setCurrentStock] = useState("");
  const [notes, setNotes] = useState("");
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleCreate = async () => {
    if (!name) {
      Alert.alert("エラー", "名前は必須です");
      return;
    }
    const stock = parseInt(currentStock, 10);
    if (isNaN(stock) || stock < 0) {
      Alert.alert("エラー", "在庫数は0以上の数値を入力してください");
      return;
    }

    setLoading(true);
    try {
      await beansApi.create({
        name,
        origin: origin || undefined,
        roast_level: roastLevel as RoastLevel,
        current_stock: stock,
        notes: notes || undefined,
      });
      router.back();
    } catch (e) {
      const msg = e instanceof ApiError ? e.message : "作成に失敗しました";
      Alert.alert("エラー", msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === "ios" ? "padding" : "height"}
    >
      <ScrollView contentContainerStyle={styles.scroll} keyboardShouldPersistTaps="handled">
        <Text style={styles.label}>名前 *</Text>
        <TextInput
          style={styles.input}
          value={name}
          onChangeText={setName}
          placeholder="例: エチオピア イルガチェフェ"
          placeholderTextColor={colors.textTertiary}
          underlineColorAndroid="transparent"
        />

        <Text style={styles.label}>産地</Text>
        <TextInput
          style={styles.input}
          value={origin}
          onChangeText={setOrigin}
          placeholder="例: エチオピア"
          placeholderTextColor={colors.textTertiary}
          underlineColorAndroid="transparent"
        />

        <Text style={styles.label}>焙煎度</Text>
        <TextInput
          style={styles.input}
          value={roastLevel}
          onChangeText={setRoastLevel}
          placeholder="例: 中煎り"
          placeholderTextColor={colors.textTertiary}
          underlineColorAndroid="transparent"
        />

        <Text style={styles.label}>在庫数 (g) *</Text>
        <TextInput
          style={styles.input}
          value={currentStock}
          onChangeText={setCurrentStock}
          placeholder="例: 200"
          keyboardType="numeric"
          placeholderTextColor={colors.textTertiary}
          underlineColorAndroid="transparent"
        />

        <Text style={styles.label}>メモ</Text>
        <TextInput
          style={[styles.input, styles.textArea]}
          value={notes}
          onChangeText={setNotes}
          placeholder="メモを入力"
          multiline
          numberOfLines={4}
          placeholderTextColor={colors.textTertiary}
          underlineColorAndroid="transparent"
        />

        <TouchableOpacity
          style={[styles.button, loading && styles.buttonDisabled]}
          onPress={handleCreate}
          disabled={loading}
        >
          <Text style={styles.buttonText}>{loading ? "追加中..." : "追加する"}</Text>
        </TouchableOpacity>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  scroll: { padding: spacing["2xl"] },
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
  button: {
    backgroundColor: colors.primary,
    borderRadius: radius.sm,
    padding: spacing.lg,
    alignItems: "center",
    marginTop: spacing["4xl"],
    ...shadows.md,
  },
  buttonDisabled: { opacity: 0.6 },
  buttonText: {
    ...typography.labelLarge,
    color: colors.textInverse,
    letterSpacing: 1,
  },
});
