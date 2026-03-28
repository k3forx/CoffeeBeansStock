import { Slot, usePathname, useRouter } from "expo-router";
import { View, Text, TouchableOpacity, StyleSheet } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { colors, typography, spacing } from "@/theme";

const tabs = [
  { name: "/(tabs)", label: "珈琲豆", icon: "☕" },
  { name: "/(tabs)/profile", label: "マイページ", icon: "👤" },
] as const;

function getHeaderConfig(pathname: string) {
  if (pathname.startsWith("/beans/") && pathname !== "/beans/create") {
    return { title: "詳細", showBack: true };
  }
  if (pathname === "/beans/create") {
    return { title: "新しい珈琲豆", showBack: true };
  }
  if (pathname === "/profile") {
    return { title: "マイページ", showBack: false };
  }
  return { title: "珈琲豆", showBack: false };
}

export default function TabLayout() {
  const pathname = usePathname();
  const router = useRouter();

  const { title, showBack } = getHeaderConfig(pathname);

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        {showBack ? (
          <TouchableOpacity style={styles.headerLeft} onPress={() => router.back()}>
            <Text style={styles.headerBackText}>‹ 戻る</Text>
          </TouchableOpacity>
        ) : (
          <View style={styles.headerLeft} />
        )}
        <Text style={styles.headerTitle} numberOfLines={1}>
          {title}
        </Text>
        <View style={styles.headerRight} />
      </View>
      <View style={styles.content}>
        <Slot />
      </View>
      <View style={styles.tabBar}>
        {tabs.map((tab) => {
          const isProfile = pathname === "/profile";
          const isFocused = tab.name === "/(tabs)/profile" ? isProfile : !isProfile;
          return (
            <TouchableOpacity
              key={tab.name}
              style={styles.tab}
              onPress={() => router.replace(tab.name)}
            >
              {isFocused && <View style={styles.tabIndicator} />}
              <Text style={[styles.tabIcon, isFocused && styles.tabIconActive]}>
                {tab.icon}
              </Text>
              <Text style={[styles.tabLabel, isFocused && styles.tabLabelActive]}>
                {tab.label}
              </Text>
            </TouchableOpacity>
          );
        })}
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  header: {
    backgroundColor: colors.background,
    flexDirection: "row",
    alignItems: "center",
    paddingVertical: 14,
    paddingHorizontal: spacing.xl,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.border,
  },
  headerLeft: { width: 70 },
  headerBackText: {
    ...typography.bodyMedium,
    color: colors.primaryMuted,
  },
  headerTitle: {
    flex: 1,
    ...typography.titleMedium,
    color: colors.textPrimary,
    textAlign: "center",
    letterSpacing: 2,
  },
  headerRight: { width: 70 },
  content: { flex: 1 },
  tabBar: {
    flexDirection: "row",
    backgroundColor: colors.surface,
    borderTopWidth: StyleSheet.hairlineWidth,
    borderTopColor: colors.border,
    paddingBottom: 4,
    paddingTop: 10,
  },
  tab: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
    position: "relative",
  },
  tabIndicator: {
    position: "absolute",
    top: -10,
    width: 16,
    height: 1.5,
    borderRadius: 1,
    backgroundColor: colors.primary,
  },
  tabIcon: { fontSize: 18, opacity: 0.3 },
  tabIconActive: { opacity: 0.9 },
  tabLabel: {
    ...typography.labelSmall,
    color: colors.textTertiary,
    marginTop: 3,
  },
  tabLabelActive: { color: colors.primary },
});
