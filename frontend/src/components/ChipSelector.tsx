import { View, Text, TouchableOpacity, StyleSheet } from "react-native";
import { colors, typography, spacing, radius } from "@/theme";

interface ChipItem<T extends string> {
  label: string;
  value: T;
}

interface ChipSelectorProps<T extends string> {
  items: ChipItem<T>[];
  selected: T | "";
  onSelect: (value: T) => void;
}

export function ChipSelector<T extends string>({ items, selected, onSelect }: ChipSelectorProps<T>) {
  return (
    <View style={styles.chipContainer}>
      {items.map((item) => (
        <TouchableOpacity
          key={item.value}
          style={[styles.chip, selected === item.value && styles.chipSelected]}
          onPress={() => onSelect(item.value)}
        >
          <Text style={[styles.chipText, selected === item.value && styles.chipTextSelected]}>
            {item.label}
          </Text>
        </TouchableOpacity>
      ))}
    </View>
  );
}

const styles = StyleSheet.create({
  chipContainer: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: spacing.sm,
  },
  chip: {
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.lg,
    borderRadius: radius.full,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.background,
  },
  chipSelected: {
    backgroundColor: colors.primary,
    borderColor: colors.primary,
  },
  chipText: {
    ...typography.bodyMedium,
    color: colors.textSecondary,
  },
  chipTextSelected: {
    color: colors.textInverse,
  },
});
