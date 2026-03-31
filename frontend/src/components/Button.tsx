import { Text, TouchableOpacity, StyleSheet, type StyleProp, type ViewStyle, type TextStyle } from "react-native";
import { colors, typography, spacing, radius, shadows } from "@/theme";
import { formStyles } from "@/theme/formStyles";

type ButtonProps = {
  title: string;
  onPress: () => void;
  variant?: "primary" | "secondary";
  size?: "default" | "small";
  disabled?: boolean;
  loading?: boolean;
  loadingText?: string;
  shadow?: boolean;
  style?: StyleProp<ViewStyle>;
  textStyle?: StyleProp<TextStyle>;
};

export function Button({
  title,
  onPress,
  variant = "primary",
  size = "default",
  disabled = false,
  loading = false,
  loadingText,
  shadow = false,
  style,
  textStyle,
}: ButtonProps) {
  const isDisabled = disabled || loading;
  const isPrimary = variant === "primary";
  const isSmall = size === "small";

  const containerStyle: ViewStyle[] = [
    isSmall ? styles.smallBase : styles.defaultBase,
    isPrimary ? styles.primaryBg : styles.secondaryBg,
    shadow && isPrimary ? shadows.md : {},
    isDisabled ? styles.disabled : {},
  ];

  const labelStyle: TextStyle[] = [
    isSmall
      ? isPrimary ? styles.smallTextPrimary : styles.smallTextSecondary
      : isPrimary ? styles.defaultTextPrimary : styles.defaultTextSecondary,
  ];

  return (
    <TouchableOpacity
      style={[...containerStyle, style]}
      onPress={onPress}
      disabled={isDisabled}
      activeOpacity={0.7}
    >
      <Text style={[...labelStyle, textStyle]}>
        {loading && loadingText ? loadingText : title}
      </Text>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  defaultBase: {
    ...formStyles.buttonBase,
  },
  smallBase: {
    borderRadius: radius.sm,
    paddingVertical: spacing.xs,
    paddingHorizontal: spacing.md,
    alignItems: "center",
  },
  primaryBg: {
    backgroundColor: colors.primary,
  },
  secondaryBg: {
    backgroundColor: "transparent",
    borderWidth: 1,
    borderColor: colors.border,
  },
  disabled: {
    opacity: 0.6,
  },
  defaultTextPrimary: {
    ...formStyles.buttonTextBase,
  },
  defaultTextSecondary: {
    ...typography.labelLarge,
    color: colors.textSecondary,
  },
  smallTextPrimary: {
    ...typography.labelSmall,
    color: colors.textInverse,
  },
  smallTextSecondary: {
    ...typography.labelSmall,
    color: colors.textSecondary,
  },
});
