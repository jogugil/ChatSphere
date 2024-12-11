import { useEffect } from 'react';

const usePolling = (fetchData: Function, setData: Function, interval: number) => {
  useEffect(() => {
    const fetchInterval = setInterval(async () => {
      const data = await fetchData();
      setData(data);
    }, interval);

    return () => clearInterval(fetchInterval);
  }, [fetchData, setData, interval]);
};

export default usePolling;

