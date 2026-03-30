import { useState } from "react";
import {
  Text,
  TouchableOpacity,
  StyleSheet,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  Alert,
} from "react-native";
import { useRouter } from "expo-router";
import { beansApi } from "../../../src/api/beans";
import { ApiError } from "../../../src/api/client";
import type { RoastLevel, RoastDetail } from "../../../src/types/api";
import { colors, typography, spacing, radius, shadows } from "@/theme";
import { ROAST_LEVELS, ROAST_DETAILS } from "../../../src/constants/roastLevels";
import { ChipSelector } from "../../../src/components/ChipSelector";
import { validateBeanForm } from "../../../src/utils/validation";
import { FormInput } from "../../../src/components/FormInput";

export default function CreateBeanScreen() {
  const [name, setName] = useState("");
  const [origin, setOrigin] = useState("");
  const [roastLevel, setRoastLevel] = useState<RoastLevel | "">("");
  const [roastDetail, setRoastDetail] = useState<RoastDetail | "">("");
  const [currentStock, setCurrentStock] = useState("");
  const [notes, setNotes] = useState("");
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const router = useRouter();

  const handleRoastLevelSelect = (level: RoastLevel) => {
    setRoastLevel(level);
    setRoastDetail("");
  };

  const handleCreate = async () => {
    const result = validateBeanForm({ name, roastLevel, currentStock });
    if (!result.valid) {
      setErrors(result.errors);
      return;
    }
    setErrors({});

    setLoading(true);
    try {
      await beansApi.create({
        name,
        origin: origin || undefined,
        roast_level: roastLevel as RoastLevel,
        roast_detail: roastDetail || undefined,
        current_stock: result.stock,
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
        <FormInput
          label="名前"
          required
          value={name}
          onChangeText={(v) => { setName(v); setErrors(prev => { const { name: _, ...rest } = prev; return rest; }); }}
          error={errors.name}
          placeholder="例: エチオピア イルガチェフェ"
        />

        <FormInput
          label="産地"
          value={origin}
          onChangeText={setOrigin}
          placeholder="例: エチオピア"
        />

        <Text style={styles.label}>焙煎度 <Text style={styles.required}>*</Text></Text>
        <ChipSelector
          items={ROAST_LEVELS}
          selected={roastLevel}
          onSelect={handleRoastLevelSelect}
        />

        {roastLevel !== "" && (
          <>
            <Text style={styles.label}>焙煎度（詳細）</Text>
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
          placeholder="例: 200"
          keyboardType="numeric"
        />

        <FormInput
          label="メモ"
          value={notes}
          onChangeText={setNotes}
          placeholder="メモを入力"
          multiline
          numberOfLines={4}
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
  required: {
    color: colors.danger,
  },
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
