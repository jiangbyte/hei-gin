/** Author: Charlie */

import 'vue-router'

declare module 'vue-router' {
  /**
   * Vue Router 路由元信息。
   *
   * 这里不再使用 title、hide、keepAlive 等旧字段，统一直接承载后端 SysResource 字段。
   * 动态路由生成、菜单展示、标签页缓存等逻辑都从这些资源字段读取。
   */
  interface RouteMeta {
    // 资源主键。
    id?: string

    // 父资源 ID，用于菜单树和隐藏页面高亮回溯。
    parent_id?: string | null

    // 资源编码。同一资源模块内唯一；动态 route.name 使用 module_id + code。
    code?: string

    // 资源名称，用于菜单、面包屑、标签页、文档标题等展示文本。
    name?: string

    // 资源类型，决定该资源是否生成路由、菜单或仅作为权限节点。
    resource_type?: AppRoute.ResourceType

    // 所属资源模块 ID，保留给后续按模块分组、筛选或权限隔离使用。
    module_id?: string | null

    // 所属资源模块名称，用于回显。
    module_id_name?: string | null

    // 前端路由路径，例如 /dashboard。
    path?: string | null

    // 页面组件路径，约定为 src/views 下的相对路径，例如 /dashboard/index.vue。
    component?: string | null

    // 路由重定向地址，目录未配置时会由前端自动补第一个可见子节点。
    redirect?: string | null

    // Iconify 图标名。
    icon?: string | null

    // 资源颜色。
    color?: string | null

    // 外链地址；存在 href 时路由守卫会打开新窗口并阻止当前页面跳转。
    href?: string | null

    // 排序值，数字越小越靠前。
    sort?: number

    // 是否在菜单和搜索中可见；隐藏页面仍可生成路由。
    is_visible?: boolean

    // 是否加入 keep-alive 缓存列表。
    is_cache?: boolean

    // 是否作为固定标签页展示。
    is_affix?: boolean

    // 资源状态，当前只有 ENABLED 资源会参与前端路由转换。
    status?: AppRoute.ResourceStatus

    // 资源描述，前端暂不展示，保留给后续扩展。
    description?: string | null

    // 布局类型，null/"default" 使用后台 Layout，"fullscreen" 为独立全屏页面。
    layout?: string | null

    // 是否作为独立全屏页面渲染，不挂载到后台 Layout 下。
    is_fullscreen?: boolean
  }
}

declare global {
  /**
   * 应用路由相关全局类型。
   *
   * 声明文件顶部存在 import，因此这里必须使用 declare global 才能让 AppRoute
   * 在项目其它文件中以全局命名空间形式访问。
   */
  namespace AppRoute {
    // 后端 ResourceType 枚举。只有 CATALOG/MENU/PAGE 会进入前端路由菜单体系。
    type ResourceType = 'CATALOG' | 'MENU' | 'PAGE' | 'BUTTON' | 'ACTION' | 'API_GROUP'

    // 后端通用状态枚举。当前资源路由只消费 ENABLED。
    type ResourceStatus = 'ENABLED' | 'DISABLED'

    /**
     * 后端 SysResource 原始资源结构。
     *
     * staticRoutes 和后续动态接口返回值都应遵循该字段结构，避免前端再维护一套旧路由字段。
     */
    interface RowRoute {
      id: string
      parent_id: string | null
      code: string
      name: string
      resource_type: ResourceType
      module_id?: string | null
      module_id_name?: string | null
      path?: string | null
      component?: string | null
      redirect?: string | null
      icon?: string | null
      color?: string | null
      href?: string | null
      sort: number
      is_visible: boolean
      is_cache: boolean
      is_affix: boolean
      status: ResourceStatus
      description?: string | null
      layout?: string | null
      created_at?: string | null
      created_by?: string | null
      updated_at?: string | null
      updated_by?: string | null
      is_fullscreen?: boolean
    }

    /**
     * 前端按 module_id 分组后的资源模块结构。
     *
     * 后端已不再按模块聚合返回资源，改为返回扁平 RowRoute[]，
     * 前端在 route store 中通过 groupResourcesByModule 自行分组。
     */
    interface ResourceModule {
      id: string
      name: string
      code: string
      client: string
      icon?: string | null
      color?: string | null
      sort: number
      resources: RowRoute[]
    }

    /**
     * 前端内部路由节点。
     *
     * 它由 RowRoute 标准化而来：route.name 使用 module_id + code，meta 保存完整资源字段。
     */
    interface Route extends Omit<RowRoute, 'name' | 'path' | 'component' | 'redirect'> {
      name: string
      path: string
      redirect?: string
      component?: import('vue-router').RouteRecordRaw['component']
      children?: Route[]
      meta: import('vue-router').RouteMeta
    }

    /**
     * 轻量菜单类型。
     *
     * 不直接使用 naive-ui 的 MenuOption，是为了避免 Pinia store state 展开过深类型，
     * 触发 类型 instantiation is excessively deep and possibly infinite。
     */
    interface MenuOption {
      // 菜单唯一键，当前使用资源 path，便于和路由路径同步高亮。
      key: string

      // 菜单显示内容；目录为纯文本，页面资源为 RouterLink 渲染函数。
      label?: string | (() => unknown)

      // 图标渲染函数。
      icon?: () => unknown

      // 子菜单。
      children?: MenuOption[]

      // 资源 ID，供 arrayToTree 组装树结构。
      id?: string

      // 父资源 ID，供 arrayToTree 组装树结构。
      parent_id?: string | null
    }
  }
}
