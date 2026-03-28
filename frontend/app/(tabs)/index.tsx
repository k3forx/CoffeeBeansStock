import { useCallback, useState } from "react";
import {
  View,
  Text,
  FlatList,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  RefreshControl,
  Alert,
} from "react-native";
import { useFocusEffect, useRouter } from "expo-router";
import { beansApi } from "../../src/api/beans";
import type { CoffeeBean, RoastLevel } from "../../src/types/api";
import { colors, typography, spacing, radius, shadows, getStockColor } from "@/theme";

const ROAST_LEVEL_LABELS: Record<RoastLevel, string> = {
  shallow: "浅煎り",
  medium: "中煎り",
  medium_deep: "中深煎り",
  deep: "深煎り",
};

export default function BeansListScreen() {
  const [beans, setBeans] = useState<CoffeeBean[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const router = useRouter();

  const fetchBeans = async () => {
    try {
      const result = await beansApi.list(100, 0);
      setBeans(result.beans);
    } catch {
      Alert.alert("エラー", "データの取得に失敗しました");
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  useFocusEffect(
    useCallback(() => {
      fetchBeans();
    }, [])
  );

  const onRefresh = () => {
    setRefreshing(true);
    fetchBeans();
  };

  const renderBean = ({ item }: { item: CoffeeBean }) => {
    const stockColor = getStockColor(item.current_stock);
    return (
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
    );
  };

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <FlatList
        data={beans}
        keyExtractor={(item) => item.id}
        renderItem={renderBean}
        contentContainerStyle={beans.length === 0 ? styles.center : styles.list}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={onRefresh}
            tintColor={colors.primary}
          />
        }
        ListEmptyComponent={
          <View style={styles.empty}>
            <Text style={styles.emptyIcon}>☕</Text>
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
        <Text style={styles.fabText}>＋</Text>
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
    width: 50,
    height: 50,
    borderRadius: 25,
    borderWidth: 1,
    justifyContent: "center",
    alignItems: "center",
    flexShrink: 0,
  },
  stockNum: {
    ...typography.stockNumber,
    lineHeight: 18,
  },
  stockUnit: {
    fontSize: 9,
    letterSpacing: 0.5,
    opacity: 0.7,
  },
  empty: { alignItems: "center" },
  emptyIcon: { fontSize: 48, marginBottom: spacing.lg, opacity: 0.4 },
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
    width: 52,
    height: 52,
    borderRadius: 26,
    backgroundColor: colors.primary,
    justifyContent: "center",
    alignItems: "center",
    ...shadows.lg,
  },
  fabText: {
    color: colors.textInverse,
    fontSize: 24,
    fontWeight: "300",
    lineHeight: 26,
  },
});
