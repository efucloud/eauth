import component from './en-US/component';
import global_page from './en-US/global_page';
import menu from './en-US/menu';
import model from './en-US/model';
import model_application from './en-US/model_application';
import model_face_recognition from './en-US/model_face_recognition';
import model_provider_oidc from './en-US/model_provider_oidc';
import model_user from './en-US/model_user';
import model_user_auth_profile from './en-US/model_user_auth_profile';
import model_user_token from './en-US/model_user_token';
import pages from './en-US/pages';
import pwa from './en-US/pwa';
import settings from './en-US/settings';

export default {
  ...pages,
  ...global_page,
  ...menu,
  ...settings,
  ...pwa,
  ...component,
  ...model,
  ...model_user,
  ...model_application,
  ...model_user_token,
  ...model_user_auth_profile,
  ...model_face_recognition,
  ...model_provider_oidc,
};
