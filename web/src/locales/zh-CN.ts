import component from './zh-CN/component';
import global_page from './zh-CN/global_page';
import menu from './zh-CN/menu';
import model from './zh-CN/model';
import model_application from './zh-CN/model_application';
import model_face_recognition from './zh-CN/model_face_recognition';
import model_provider_oidc from './zh-CN/model_provider_oidc';
import model_user from './zh-CN/model_user';
import model_user_auth_profile from './zh-CN/model_user_auth_profile';
import model_user_token from './zh-CN/model_user_token';
import pages from './zh-CN/pages';
import pwa from './zh-CN/pwa';
import settings from './zh-CN/settings';
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
