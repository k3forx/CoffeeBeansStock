import { useState } from "react";
import {
  Text,
  TouchableOpacity,
  StyleSheet,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
} from "react-native";
import { useRouter } from "expo-router";
import { beansApi } from "../../../src/api/beans";
import { showApiError } from "../../../src/utils/errorHandler";
import type { RoastLevel, RoastDetail } from "../../../src/types/api";
import { colors, spacing, shadows, formStyles } from "@/theme";
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
      showApiError(e, "作成に失敗しました");
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
          style={[styles.button, loading && formStyles.buttonDisabled]}
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
  button: {
    ...formStyles.buttonBase,
    backgroundColor: colors.primary,
    marginTop: spacing["4xl"],
    ...shadows.md,
  },
  buttonText: {
    ...formStyles.buttonTextBase,
    letterSpacing: 1,
  },
});
