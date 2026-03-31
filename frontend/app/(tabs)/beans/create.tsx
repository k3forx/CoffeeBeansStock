import { useState } from "react";
import {
  Text,
  StyleSheet,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
} from "react-native";
import { useRouter } from "expo-router";
import { beansApi } from "../../../src/api/beans";
import { showApiError } from "../../../src/utils/errorHandler";
import { colors, spacing, formStyles } from "@/theme";
import { Button } from "@/components/Button";
import { ROAST_LEVELS, ROAST_DETAILS } from "../../../src/constants/roastLevels";
import { ChipSelector } from "../../../src/components/ChipSelector";
import { FormInput } from "../../../src/components/FormInput";
import { useBeanForm } from "../../../src/hooks/useBeanForm";

export default function CreateBeanScreen() {
  const form = useBeanForm();
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleCreate = async () => {
    const result = form.validate();
    if (!result.valid) {
      return;
    }

    setLoading(true);
    try {
      await beansApi.create(form.toCreateInput(result.stock));
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
          value={form.fields.name}
          onChangeText={(v) => form.setField("name", v)}
          error={form.errors.name}
          placeholder="例: エチオピア イルガチェフェ"
        />

        <FormInput
          label="産地"
          value={form.fields.origin}
          onChangeText={(v) => form.setField("origin", v)}
          placeholder="例: エチオピア"
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
          placeholder="例: 200"
          keyboardType="numeric"
        />

        <FormInput
          label="メモ"
          value={form.fields.notes}
          onChangeText={(v) => form.setField("notes", v)}
          placeholder="メモを入力"
          multiline
          numberOfLines={4}
        />

        <Button
          title="追加する"
          onPress={handleCreate}
          loading={loading}
          loadingText="追加中..."
          shadow
          style={{ marginTop: spacing["4xl"] }}
          textStyle={{ letterSpacing: 1 }}
        />
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  scroll: { padding: spacing["2xl"] },
});
