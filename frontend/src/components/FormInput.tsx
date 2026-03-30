import React, { useState } from "react";
import {
  View,
  TextInput,
  Text,
  StyleSheet,
  KeyboardTypeOptions,
} from "react-native";
import { colors, typography, spacing } from "@/theme";

type FormInputProps = {
  label: string;
  required?: boolean;
  value: string;
  onChangeText: (text: string) => void;
  error?: string;
  placeholder?: string;
  keyboardType?: KeyboardTypeOptions;
  multiline?: boolean;
  numberOfLines?: number;
};

export function FormInput({
  label,
  required,
  value,
  onChangeText,
  error,
  placeholder,
  keyboardType,
  multiline,
  numberOfLines,
}: FormInputProps) {
  const [focused, setFocused] = useState(false);

  const borderBottomColor = error
    ? colors.danger
    : focused
      ? colors.primary
      : colors.border;
  const borderBottomWidth = focused ? 2 : 1;

  return (
    <View style={styles.container}>
      <Text style={styles.label}>
        {label}
        {required && <Text style={styles.required}> *</Text>}
      </Text>
      <TextInput
        style={[styles.input, { borderBottomColor, borderBottomWidth }]}
        value={value}
        onChangeText={onChangeText}
        placeholder={placeholder}
        placeholderTextColor={colors.textTertiary}
        keyboardType={keyboardType}
        multiline={multiline}
        numberOfLines={numberOfLines}
        onFocus={() => setFocused(true)}
        onBlur={() => setFocused(false)}
      />
      {error ? <Text style={styles.error}>{error}</Text> : null}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginBottom: spacing.md,
  },
  label: {
    ...typography.labelMedium,
    color: colors.textPrimary,
    marginBottom: spacing.xs,
  },
  required: {
    color: colors.danger,
  },
  input: {
    ...typography.bodyLarge,
    color: colors.textPrimary,
    backgroundColor: "transparent",
    borderWidth: 0,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
    paddingTop: spacing.md,
    paddingBottom: spacing.xl,
  },
  error: {
    ...typography.bodySmall,
    color: colors.danger,
    marginTop: spacing.xs,
  },
});
