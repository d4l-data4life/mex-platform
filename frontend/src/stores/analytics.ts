import { Persistor } from './utils/persistor';
import createPersistedStore from './utils/persisted-store';
import { ANALYTICS_IS_DNT } from 'config';
import { generateVisitorId } from 'utils/analytics';

export enum AnalyticsConsent {
  TRACKING = 'tracking',
}

interface StateType {
  consents?: AnalyticsConsent[];
  visitorId?: string;
}

const localPersistor = new Persistor('localStorage' in window ? localStorage : null);
const store = createPersistedStore<StateType>(localPersistor, 'analytics', {});

class AnalyticsStore {
  public onConsentsChange = store.onChange.bind(this, 'consents');

  get consents() {
    return store.get('consents') ?? [];
  }

  set consents(consents: AnalyticsConsent[]) {
    store.set('consents', consents);
  }

  get hasChosen() {
    return !!store.get('consents') || ANALYTICS_IS_DNT();
  }

  get visitorId() {
    const visitorId = store.get('visitorId');
    !visitorId && store.set('visitorId', generateVisitorId());
    return visitorId ?? store.get('visitorId');
  }

  reset(): void {
    store.reset();
  }
}

export default new AnalyticsStore();
