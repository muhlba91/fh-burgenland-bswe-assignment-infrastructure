import * as harbor from '@pulumiverse/harbor';

/**
 * Configures Harbor authentication using Dex.
 */
export const configureHarborAuth = () => {
  new harbor.ConfigAuth('harbor-auth-dex', {
    authMode: 'oidc_auth',
    oidcAutoOnboard: true,
    oidcClientId: 'harbor',
    oidcClientSecret: process.env.DEX_HARBOR_CLIENT_SECRET,
    oidcEndpoint: 'https://auth.hochschule-burgenland.muehlbachler.xyz/api/dex',
    oidcName: 'GITHUB',
    oidcUserClaim: 'preferred_username',
    oidcGroupsClaim: 'groups',
    oidcScope: 'openid,email,profile,groups,offline_access',
  });
};
