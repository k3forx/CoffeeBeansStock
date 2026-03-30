import React, { useEffect, useRef } from "react";
import { Animated, StyleSheet, View } from "react-native";
import { colors, spacing, radius, shadows } from "@/theme";

function BeanCardSkeleton() {
  const opacity = useRef(new Animated.Value(0.3)).current;

  useEffect(() => {
    const animation = Animated.loop(
      Animated.sequence([
        Animated.timing(opacity, {
          toValue: 0.7,
          duration: 800,
          useNativeDriver: true,
        }),
        Animated.timing(opacity, {
          toValue: 0.3,
          duration: 800,
          useNativeDriver: true,
        }),
      ])
    );
    animation.start();
    return () => animation.stop();
  }, [opacity]);

  return (
    <Animated.View style={[styles.card, { opacity }]}>
      <View style={styles.cardContent}>
        <View style={styles.left}>
          <View style={[styles.line, styles.lineWide]} />
          <View style={[styles.line, styles.lineNarrow]} />
        </View>
        <View style={styles.circle} />
      </View>
    </Animated.View>
  );
}

export function BeanListSkeleton() {
  return (
    <View style={styles.list}>
      <BeanCardSkeleton />
      <BeanCardSkeleton />
      <BeanCardSkeleton />
      <BeanCardSkeleton />
    </View>
  );
}

const styles = StyleSheet.create({
  list: {
    padding: spacing.lg,
    gap: spacing.md,
  },
  card: {
    backgroundColor: colors.surface,
    borderRadius: radius.md,
    ...shadows.sm,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: colors.border,
    padding: spacing.lg,
  },
  cardContent: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
  },
  left: {
    gap: spacing.sm,
  },
  line: {
    height: 14,
    borderRadius: radius.sm,
    backgroundColor: colors.border,
  },
  lineWide: {
    width: 160,
  },
  lineNarrow: {
    width: 100,
  },
  circle: {
    width: 58,
    height: 58,
    borderRadius: 29,
    backgroundColor: colors.border,
  },
});
