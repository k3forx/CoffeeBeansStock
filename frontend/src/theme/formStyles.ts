import { colors, typography, spacing, radius } from "./index";

export const formStyles = {
  label: {
    ...typography.labelMedium,
    color: colors.textSecondary,
    marginBottom: spacing.sm,
    marginTop: spacing.xl,
  },
  required: {
    color: colors.danger,
  },
  buttonDisabled: {
    opacity: 0.6,
  },
  buttonBase: {
    borderRadius: radius.sm,
    padding: spacing.lg,
    alignItems: "center" as const,
  },
  buttonTextBase: {
    ...typography.labelLarge,
    color: colors.textInverse,
  },
} as const;
