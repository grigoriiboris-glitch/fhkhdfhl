<template>
  <div class="p-4">
    <h1 class="page-title">User - <span class="fw-semi-bold">Profile</span></h1>
    <div class="row">
      <div class="col-lg-6 col-12">
        <Widget>
          <div class="widget-top-overflow text-white">
            <div class="height-250 overflow-hidden">
              <img class="img-fluid" src="/images/pictures/10.jpg" alt="picture" />
            </div>
            <a href="#" class="btn btn-outline btn-sm mb-2">
              <i class="fa fa-twitter me-2"></i>
              Follow
            </a>
          </div>
          <div class="row">
            <div class="col-md-5 col-12 text-center">
              <div class="profileContactContainer">
                <span class="thumb-xl mb-3">
                  <img class="profileAvatar rounded-circle" src="/images/pictures/10.jpg" alt="man" />
                </span>
                <h5 class="fw-normal">Adam <span class="fw-semi-bold">Johns</span></h5>
                <p>UI/UX designer</p>
                <a href="#" class="btn btn-success btn-sm mb-3">
                  &nbsp;Send <i class="fa fa-envelope ms-2"></i>&nbsp;
                </a>
                <div>
                  <ul class="profileContacts mt-sm">
                    <li><i class="fa fa-lg fa-phone fa-fw me-2"></i><a href="#"> +375 29 555-55-55</a></li>
                    <li><i class="fa fa-lg fa-envelope fa-fw me-2"></i><a href="#"> psmith@example.com</a></li>
                    <li><i class="fa fa-lg fa-map-marker fa-fw me-2"></i><a href="#"> Minsk, Belarus</a></li>
                  </ul>
                </div>
              </div>
            </div>
            <div class="col-md-7 col-12">
              <div class="stats-row mt-3">
                <div class="profileStat stat-item">
                  <p class="profileStatValue value text-end">251</p>
                  <h6 class="name">Posts</h6>
                </div>
                <div class="profileStat stat-item">
                  <p class="profileStatValue value text-end">9.38%</p>
                  <h6 class="name">Conversion</h6>
                </div>
                <div class="profileStat stat-item">
                  <p class="profileStatValue value text-end">842</p>
                  <h6 class="name">Followers</h6>
                </div>
              </div>
              <p>
                <span class="badge bg-info rounded-pill"> UI/UX </span>
                <span class="badge bg-primary rounded-pill ms-1"> Web Design </span>
                <span class="badge bg-default rounded-pill ms-1"> Mobile Apps </span>
              </p>
              <p class="lead mt-xlg">
                My name is Adam Johns and here is my new Sing user profile page.
              </p>
              <p class="text-muted">
                I love reading people's summaries page especially those who are in the same industry as me.
                Sometimes it's much easier to find your concentration during the night.
              </p>
            </div>
          </div>
        </Widget>
      </div>
      <div class="col-lg-6 col-12">
        <section class="activities">
          <h2 class="ms-3">Activities</h2>
          <div class="event">
            <header>
              <span class="eventAvatar">
                <img class="rounded-circle" src="/images/people/a5.jpg" alt="man" />
              </span>
              <h5 class="eventTitle">
                <a href="#">Bob Nilson</a> <small class="ms-2"><a href="#">@nils</a></small>
              </h5>
              <p class="eventTimestamp">February 22, 2014 at 01:59 PM</p>
            </header>
            <div class="eventBody">
              There is no such thing as maturity. There is instead an ever-evolving process of maturing.
              Because when there is a maturity, there is ...
            </div>
            <footer class="eventFooter">
              <div class="clearfix">
                <ul class="post-links mt-sm pull-left">
                  <li><a href="#"><button class="btn-link">1 hour</button></a></li>
                  <li><a href="#">
                    <span class="text-danger"><i class="fa fa-heart-o" /> Like</span>
                  </a></li>
                  <li><a href="#"><button class="btn-link">Comment</button></a></li>
                </ul>

                <span class="thumb thumb-sm pull-right">
                  <a href="#">
                    <img class="rounded-circle" alt="woman" src="/images/people/a1.jpg" />
                  </a>
                </span>
                <span class="thumb thumb-sm pull-right">
                  <a href="#"><img class="rounded-circle"
                    alt="man" src="/images/people/a5.jpg" /></a>
                </span>
                <span class="thumb thumb-sm pull-right">
                  <a href="#"><img class="rounded-circle"
                    alt="woman" src="/images/people/a3.jpg" /></a>
                </span>
              </div>
              <ul class="post-comments mt-sm">
                <li>
                  <span class="thumb-xs avatar pull-left me-sm-1">
                    <img class="rounded-circle" src="/images/people/a1.jpg" alt="woman" />
                  </span>
                  <div class="comment-body">
                    <h6 class="author fs-sm fw-semi-bold">
                      Ignacio Abad <small>6 mins ago</small>
                    </h6>
                    <p>Hey, have you heard anything about that?</p>
                  </div>
                </li>
                <li>
                  <span class="thumb-xs avatar pull-left me-sm-1">
                    <img class="rounded-circle" src="/images/avatar.png" alt="avatar-photo" />
                  </span>
                  <div class="comment-body">
                    <b-form-input class="form-control form-control-sm"
                     type="text" placeholder="Write your comment..." />
                  </div>
                </li>
              </ul>
            </footer>
          </div>
          <form class="mt" action="#">
            <div class="mb-2">
              <label for="new-event" class="sr-only">New event</label>
              <textarea id="new-event" class="form-control" placeholder="Post something..." rows="3"></textarea>
            </div>
            <div class="btn-toolbar">
              <div class="btn-group">
                <button class="btn btn-sm btn-gray-default"><i class="fa fa-camera fa-lg"></i></button>
                <button class="btn btn-sm btn-gray-default"><i class="fa fa-map-marker fa-lg"></i></button>
              </div>
              <button type="submit" class="btn btn-success text-white ms-auto">Post</button>
            </div>
          </form>
        </section>
      </div>
    </div>
  </div>
</template>

<script>
import Widget from '../components/Widget/Widget';

export default {
  name: 'Profile',
  components: { Widget },
};
</script>

<style lang="scss" scoped >
@import '../styles/app';

.profileContactContainer {
  margin-top: -75px;
}

.profileContacts {
  @include list-unstyled();

  display: inline-block;
  text-align: left;

  > li {
    margin-bottom: $spacer / 2;
  }

  > li > a {
    color: lighten($text-color, 30%);
    text-decoration: none;

    @include hover-focus {
      color: $text-color;
    }
  }
}

.profileAvatar {
  border: 3px solid $white;
}

.profileStat {
  border-left: none !important;
  padding-right: 0 !important;
}

.profileStatValue {
  font-size: 28px;
  font-weight: $font-weight-base !important;
  margin-bottom: 0;
}

.event {
  background: $white;
  border-radius: $border-radius;
  padding: 20px 20px 0;
  position: relative;
  margin-bottom: $spacer;
  box-shadow: var(--widget-shadow);
}

.eventTitle {
  margin-bottom: 2px;
  font-weight: $font-weight-semi-bold;

  a {
    text-decoration: none;
    color: #6FB0F9;
  }

  small > a {
    color: $text-muted;
  }
}

.eventAvatar {
  float: left;
  margin-right: $spacer;
}

.eventAvatar > img {
  width: 34px;
}

.eventBody {
  font-size: 0.9rem;
  margin-bottom: $spacer;
}

.eventFooter {
  background-color: $gray-100;
  margin: 20px -20px 0;
  padding: 10px 20px;
}

.eventTimestamp {
  color: $text-muted;
}

</style>
