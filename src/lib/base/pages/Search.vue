<template>
  <div class="search-page">
    <h1 class="page-title">Matching - <span class="fw-semi-bold">Results</span></h1>
    <div class="btn-toolbar justify-content-between">
      <div class="d-inline-flex">
        <div class="dropdown b-dropdown btn-group" id="__BVID__369"><button aria-haspopup="menu" aria-expanded="false"
            type="button" class="btn dropdown-toggle btn-gray-default" id="__BVID__369__BV_toggle_">Popular</button>
          <ul role="menu" tabindex="-1" class="dropdown-menu" aria-labelledby="__BVID__369__BV_toggle_">
            <li role="presentation"><a role="menuitem" href="#" target="_self" class="dropdown-item">All</a></li>
            <li role="presentation"><a role="menuitem" href="#" target="_self" class="dropdown-item">Popular</a></li>
            <li role="presentation"><a role="menuitem" href="#" target="_self" class="dropdown-item">Interesting</a>
            </li>
            <li role="presentation"><a role="menuitem" href="#" target="_self" class="dropdown-item">Latest</a></li>
          </ul>
        </div>
        <div class="dropdown b-dropdown btn-group" id="__BVID__378"><button aria-haspopup="menu" aria-expanded="false"
            type="button" class="btn dropdown-toggle btn-gray-default" id="__BVID__378__BV_toggle_">All Time</button>
          <ul role="menu" tabindex="-1" class="dropdown-menu" aria-labelledby="__BVID__378__BV_toggle_">
            <li role="presentation"><a role="menuitem" href="#" target="_self" class="dropdown-item">Last 24h</a></li>
            <li role="presentation"><a role="menuitem" href="#" target="_self" class="dropdown-item">Last Month</a></li>
            <li role="presentation"><a role="menuitem" href="#" target="_self" class="dropdown-item">Last Month</a></li>
          </ul>
        </div>
      </div>
      <div role="group" class="btn-group"><button type="button" class="btn active btn-gray"><i
            class="fa fa-th-list"></i></button><button type="button" class="btn btn-gray"><i
            class="fa fa-th-large"></i></button></div>
    </div>
    <div class="row mt-3 d-block">
      <div class="float-lg-end col-sm-12 col-lg-3">
        <h5>Results <span class="fw-semi-bold">Filtering</span></h5>
        <p class="text-muted fs-mini">Listed content is categorized by the following groups:</p>
        <ul class="nav flex-column nav-stacked searchResultCategories mt nav-pills">
          <li class="nav-item"><a href="#" target="_self" class="nav-link">
              Hot Ideas
              <span class="badge bg-danger rounded-pill float-end">34</span></a></li>
          <li class="nav-item"><a href="#" target="_self" class="nav-link">
              Latest Pictures
              <span class="badge bg-success rounded-pill float-end">9</span></a></li>
          <li class="nav-item"><a href="#" target="_self" class="nav-link">Labels of Day</a></li>
          <li class="nav-item"><a href="#" target="_self" class="nav-link">Recent Movies</a></li>
          <li class="nav-item"><a href="#" target="_self" class="nav-link">
              Globals
              <span class="badge bg-info rounded-pill float-end">18</span></a></li>
        </ul>
      </div>
      <div class="col-sm-12 col-lg-9">
        <p class="searchResultsCount">About 94 700 000 (0.39 sec.) results</p>
        <section class="searchResultItem"><a href="#" class="imageLink"><img src="/images/pictures/1.jpg"
              alt="mountains" class="image"></a>
          <div class="searchResultItemBody">
            <div class="row">
              <div class="col-md-9">
                <h4 class="searchResultItemHeading"><a href="#">Next generation admin template</a></h4>
                <p class="info">New York, NY 20188</p>
                <p class="description">
                  Not just usual Metro. But something bigger. Not just usual widgets, but real
                  widgets. Not just yet another admin template,
                  but next generation admin template.
                </p>
              </div>
              <div xs="12" class="text-center col-md-3">
                <p class="value3 mt-sm">$9, 700</p>
                <p class="fs-mini text-muted">PER WEEK</p>
                <button type="button" class="btn btn-info btn-sm">Learn More</button>
              </div>
            </div>
          </div>
        </section>
        <section class="searchResultItem"><a href="#" class="imageLink"><img src="/images/pictures/2.jpg"
              alt="mountains" class="image"></a>
          <div class="searchResultItemBody">
            <div class="row">
              <div class="col-md-9">
                <h4 class="searchResultItemHeading"><a href="#">Next generation admin template<small data-v-1ee88924=""><span data-v-1ee88924=""
                  class="badge bg-pill bg-danger float-end"><span data-v-1ee88924="" class="fw-normal"> Best
                    Deal!</span></span></small></a></h4>
                <p class="info">New York, NY 20188</p>
                <p class="description">
                  Not just usual Metro. But something bigger. Not just usual widgets, but real
                  widgets. Not just yet another admin template,
                  but next generation admin template.
                </p>
              </div>
              <div xs="12" class="text-center col-md-3">
                <p class="value3 mt-sm">$9, 700</p>
                <p class="fs-mini text-muted">PER WEEK</p>
                <button type="button" class="btn btn-info btn-sm">Learn More</button>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>


</template>

<script>
export default {
  name: 'SearchPage',
};
</script>

<style  lang="scss" scoped >
@import '../styles/app';

.search-page {
  .searchResultCategories {
    > li > a {
      color: $gray-600;
      font-weight: $font-weight-normal;

      &:hover {
        color: $gray-700;
        background-color: $gray-400;
      }
    }
  }

  .searchResultsCount {
    margin-top: 10px;
  }

  .searchResultItem {
    padding: 20px;
    background-color: $white;
    border-radius: $border-radius;
    box-shadow: var(--widget-shadow);

    &:first-of-type {
      overflow: hidden;
    }

    @include clearfix();

    .imageLink {
      display: block;
      overflow: hidden;
      border-top-left-radius: $border-radius;
      border-bottom-left-radius: $border-radius;

      @include media-breakpoint-up(md) {
        display: inline-block;
        margin: -20px 0 -20px -20px;
        float: left;
        width: 200px;
      }

      @include media-breakpoint-down(sm) {
        max-height: 200px;
      }
    }

    .image {
      max-width: 100%;
    }

    .info {
      margin-top: 2px;
      font-size: $font-size-sm;
      color: $text-muted;
    }

    .description {
      font-size: $font-size-mini;
    }

    + .searchResultItem {
      margin-top: 20px;
    }

    .searchResultItemBody {
      height: auto;

      @include media-breakpoint-down(sm) {
        margin-top: 10px;
      }

      @include media-breakpoint-up(md) {
        margin-left: 200px;
      }
    }

    .searchResultItemHeading {
      font-weight: $font-weight-normal;

      > a {
        color: $link-color;
        text-decoration: none;
      }

      @include media-breakpoint-up(md) {
        margin: 0;
      }
    }
  }
}

</style>
