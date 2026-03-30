import { useCallback, useRef } from "react";
import {
  View,
  Text,
  FlatList,
  TouchableOpacity,
  StyleSheet,
  RefreshControl,
  Animated,
} from "react-native";
import { useFocusEffect, useRouter } from "expo-router";
import { MaterialCommunityIcons, Ionicons } from "@expo/vector-icons";
import type { CoffeeBean } from "../../src/types/api";
import { colors, typography, spacing, radius, shadows, getStockColor } from "@/theme";
import { ROAST_LEVEL_LABELS } from "../../src/constants/roastLevels";
import { BeanListSkeleton } from "../../src/components/SkeletonLoader";
import { useBeansList } from "@/hooks/useBeansList";

export default function BeansListScreen() {
  const beansList = useBeansList();
  const router = useRouter();

  const animatedValues = useRef<Animated.Value[]>([]).current;

  const getAnimatedValue = (index: number) => {
    if (!animatedValues[index]) {
      animatedValues[index] = new Animated.Value(0);
      Animated.timing(animatedValues[index], {
        toValue: 1,
        duration: 300,
        delay: index * 50,
        useNativeDriver: true,
      }).start();
    }
    return animatedValues[index];
  };

  useFocusEffect(
    useCallback(() => {
      animatedValues.length = 0;
      beansList.fetchBeans();
    }, [])
  );

  const renderBean = ({ item, index }: { item: CoffeeBean; index: number }) => {
    const stockColor = getStockColor(item.current_stock);
    const anim = getAnimatedValue(index);
    return (
      <Animated.View style={{ opacity: anim, transform: [{ translateY: anim.interpolate({ inputRange: [0, 1], outputRange: [8, 0] }) }] }}>
        <TouchableOpacity
          style={styles.card}
          onPress={() => router.push(`/beans/${item.id}`)}
          activeOpacity={0.7}
        >
          <View style={[styles.accentBar, { backgroundColor: stockColor }]} />
          <View style={styles.cardContent}>
            <View style={styles.cardLeft}>
              <Text style={styles.beanName}>{item.name}</Text>
              <View style={styles.metaRow}>
                {item.origin && <Text style={styles.beanMeta}>{item.origin}</Text>}
                {item.roast_level && <Text style={styles.beanMeta}>{ROAST_LEVEL_LABELS[item.roast_level]}</Text>}
              </View>
            </View>
            <View style={[styles.stockCircle, { borderColor: stockColor }]}>
              <Text style={[styles.stockNum, { color: stockColor }]}>
                {item.current_stock}
              </Text>
              <Text style={[styles.stockUnit, { color: stockColor }]}>gram</Text>
            </View>
          </View>
        </TouchableOpacity>
      </Animated.View>
    );
  };

  if (beansList.loading) {
    return (
      <View style={styles.container}>
        <BeanListSkeleton />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <FlatList
        data={beansList.beans}
        keyExtractor={(item) => item.id}
        renderItem={renderBean}
        contentContainerStyle={beansList.beans.length === 0 ? styles.center : styles.list}
        refreshControl={
          <RefreshControl
            refreshing={beansList.refreshing}
            onRefresh={beansList.onRefresh}
            tintColor={colors.primary}
          />
        }
        ListEmptyComponent={
          <View style={styles.empty}>
            <MaterialCommunityIcons name="coffee-outline" size={48} color={colors.textTertiary} style={{ opacity: 0.4, marginBottom: spacing.lg }} />
            <Text style={styles.emptyText}>珈琲豆がまだ登録されていません</Text>
            <Text style={styles.emptySubtext}>右下の＋ボタンから追加しましょう</Text>
          </View>
        }
      />
      <TouchableOpacity
        style={styles.fab}
        onPress={() => router.push("/beans/create")}
        activeOpacity={0.8}
      >
        <Ionicons name="add" size={26} color={colors.textInverse} />
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  center: { flex: 1, justifyContent: "center", alignItems: "center" },
  list: { padding: spacing.lg },
  card: {
    backgroundColor: colors.surface,
    borderRadius: radius.md,
    padding: spacing.xl,
    paddingLeft: spacing["2xl"],
    marginBottom: spacing.md,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: colors.border,
    position: "relative",
    overflow: "hidden",
    ...shadows.sm,
  },
  accentBar: {
    position: "absolute",
    top: 0,
    left: 0,
    bottom: 0,
    width: 4,
    borderTopLeftRadius: radius.md,
    borderBottomLeftRadius: radius.md,
  },
  cardContent: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
  },
  cardLeft: { flex: 1, marginRight: spacing.md },
  beanName: {
    ...typography.titleMedium,
    color: colors.textPrimary,
    marginBottom: spacing.xs,
  },
  metaRow: { flexDirection: "row", gap: spacing.md },
  beanMeta: {
    ...typography.bodySmall,
    color: colors.textSecondary,
  },
  stockCircle: {
    width: 58,
    height: 58,
    borderRadius: 29,
    borderWidth: 1.5,
    justifyContent: "center",
    alignItems: "center",
    flexShrink: 0,
  },
  stockNum: {
    ...typography.stockNumber,
    fontSize: 17,
    lineHeight: 18,
  },
  stockUnit: {
    fontSize: 9,
    letterSpacing: 0.5,
    opacity: 0.7,
  },
  empty: { alignItems: "center" },
  emptyText: {
    ...typography.bodyLarge,
    color: colors.textSecondary,
    marginBottom: spacing.xs,
  },
  emptySubtext: {
    ...typography.bodyMedium,
    color: colors.textTertiary,
  },
  fab: {
    position: "absolute",
    right: spacing.xl,
    bottom: spacing.xl,
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: colors.primary,
    justifyContent: "center",
    alignItems: "center",
    ...shadows.lg,
  },
});
