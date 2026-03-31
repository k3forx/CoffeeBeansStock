import { useState } from "react";
import { Alert } from "react-native";
import { beansApi } from "@/api/beans";
import { usagesApi } from "@/api/usages";
import { showApiError } from "@/utils/errorHandler";
import type { CoffeeBean, UsageHistory } from "@/types/api";

interface UseUsageTrackingConfig {
  beanId: string;
  onBeanUpdated: (bean: CoffeeBean) => void;
}

export interface UseUsageTrackingReturn {
  usages: UsageHistory[];
  usageLoading: boolean;
  quickButtonLoading: number | null;
  manualGrams: string;
  manualLoading: boolean;
  deletingUsageId: string | null;
  loadUsages: () => Promise<void>;
  handleQuickUsage: (grams: number) => Promise<void>;
  handleManualUsage: () => Promise<void>;
  handleDeleteUsage: (usageId: string) => Promise<void>;
  setManualGrams: (value: string) => void;
}

function getTodayDate(): string {
  const d = new Date();
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

export function useUsageTracking({
  beanId,
  onBeanUpdated,
}: UseUsageTrackingConfig): UseUsageTrackingReturn {
  const [usages, setUsages] = useState<UsageHistory[]>([]);
  const [usageLoading, setUsageLoading] = useState(false);
  const [quickButtonLoading, setQuickButtonLoading] = useState<number | null>(null);
  const [manualGrams, setManualGrams] = useState("");
  const [manualLoading, setManualLoading] = useState(false);
  const [deletingUsageId, setDeletingUsageId] = useState<string | null>(null);

  const loadUsages = async () => {
    setUsageLoading(true);
    try {
      const data = await usagesApi.list(beanId, 10, 0);
      setUsages(data.usages);
    } catch (e) {
      showApiError(e, "使用履歴の取得に失敗しました");
    } finally {
      setUsageLoading(false);
    }
  };

  const refreshAfterUsage = async () => {
    const [updatedBean, updatedUsages] = await Promise.all([
      beansApi.get(beanId),
      usagesApi.list(beanId, 10, 0),
    ]);
    onBeanUpdated(updatedBean);
    setUsages(updatedUsages.usages);
  };

  const handleQuickUsage = async (grams: number) => {
    setQuickButtonLoading(grams);
    try {
      await usagesApi.create(beanId, {
        usage_date: getTodayDate(),
        quantity: grams,
      });
      await refreshAfterUsage();
    } catch (e) {
      showApiError(e, "記録に失敗しました", {
        INSUFFICIENT_STOCK: "在庫が不足しています",
      });
    } finally {
      setQuickButtonLoading(null);
    }
  };

  const handleManualUsage = async () => {
    const grams = parseInt(manualGrams, 10);
    if (isNaN(grams) || grams <= 0) {
      Alert.alert("エラー", "1g以上の数値を入力してください");
      return;
    }
    setManualLoading(true);
    try {
      await usagesApi.create(beanId, {
        usage_date: getTodayDate(),
        quantity: grams,
      });
      await refreshAfterUsage();
      setManualGrams("");
    } catch (e) {
      showApiError(e, "記録に失敗しました", {
        INSUFFICIENT_STOCK: "在庫が不足しています",
      });
    } finally {
      setManualLoading(false);
    }
  };

  const handleDeleteUsage = async (usageId: string) => {
    setDeletingUsageId(usageId);
    try {
      await usagesApi.delete(beanId, usageId);
      await refreshAfterUsage();
    } catch (e) {
      showApiError(e, "削除に失敗しました");
    } finally {
      setDeletingUsageId(null);
    }
  };

  return {
    usages,
    usageLoading,
    quickButtonLoading,
    manualGrams,
    manualLoading,
    deletingUsageId,
    loadUsages,
    handleQuickUsage,
    handleManualUsage,
    handleDeleteUsage,
    setManualGrams,
  };
}
