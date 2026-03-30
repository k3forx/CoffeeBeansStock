import { useState } from "react";
import { beansApi } from "@/api/beans";
import { showApiError } from "@/utils/errorHandler";
import type { CoffeeBean } from "@/types/api";

export interface UseBeansListReturn {
  beans: CoffeeBean[];
  loading: boolean;
  refreshing: boolean;
  fetchBeans: () => Promise<void>;
  onRefresh: () => void;
}

export function useBeansList(): UseBeansListReturn {
  const [beans, setBeans] = useState<CoffeeBean[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const fetchBeans = async () => {
    try {
      const result = await beansApi.list(100, 0);
      setBeans(result.beans);
    } catch (e) {
      showApiError(e, "データの取得に失敗しました");
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  const onRefresh = () => {
    setRefreshing(true);
    fetchBeans();
  };

  return { beans, loading, refreshing, fetchBeans, onRefresh };
}
